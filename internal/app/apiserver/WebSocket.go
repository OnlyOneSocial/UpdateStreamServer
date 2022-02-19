package apiserver

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/gobwas/ws"
)

//Packet ...
type Packet struct {
}

//Channel WebSocket
type Channel struct {
	*server
	conn   net.Conn
	userID int
	out    chan []byte
}

func (c *Channel) readFromStream() (string, []byte, error) {
	header, err := ws.ReadHeader(c.conn)
	if err != nil {
		return "", []byte{}, err
	}

	payload := make([]byte, header.Length)
	_, err = io.ReadFull(c.conn, payload)
	if err != nil {
		return "", []byte{}, err
	}
	if header.Masked {
		ws.Cipher(payload, header.Mask, 0)
	}

	type SocketPocket struct {
		Type string `json:"type"`
	}

	socketPocket := SocketPocket{}

	err = json.Unmarshal(payload, &socketPocket)
	if err != nil {
		fmt.Println(err)
	}

	return socketPocket.Type, payload, nil
}

func (c *Channel) writeToStream(packet []Notification) error {
	payload, err := json.Marshal(packet)
	if err != nil {
		return err
	}

	if err := ws.WriteHeader(c.conn, ws.Header{
		Fin:    true,
		Length: int64(binary.Size(payload)),
		OpCode: ws.OpText,
	}); err != nil {
		return err
	}
	//fmt.Print(string(payload[:]))
	if _, err := c.conn.Write(payload); err != nil {
		return err
	}
	return nil
}

//NewChannel init
func (c *Channel) reader() {
	//buf := bufio.NewReader(c.conn)

	for {
		typeOf, payload, err := c.readFromStream()
		if typeOf == "init" {
			type SocketPocket struct {
				Token string `json:"data"`
			}

			socketPocket := SocketPocket{}

			err = json.Unmarshal(payload, &socketPocket)
			if err != nil {
				fmt.Println(err)
			}

			userid, err := c.server.GetDataFromToken(socketPocket.Token)

			if err != nil {
				fmt.Println(err)
			}

			if _, ok := c.server.BufferNotiff[userid]; !ok {
				/*
					UserNotifications := []Notification{}
						UserNotifications = append(UserNotifications, Notification{
							Type:      "request_send",
							UserID:    10,
							UserName:  "Кейтлин",
							TimeStamp: 1645075971,
						})
						UserNotifications = append(UserNotifications, Notification{
							Type:      "request_accept",
							UserID:    10,
							UserName:  "mlebd",
							TimeStamp: 1645075971,
						})
						UserNotifications = append(UserNotifications, Notification{
							Type:      "message_send",
							UserID:    10,
							UserName:  "Astolfo",
							TimeStamp: 1645075971,
						})

						c.server.BufferNotiff[userid] = &NotiffList{List: UserNotifications}
				*/
			}

			c.userID = userid
		}

		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func (c *Channel) writer() {

	for {
		if notifications, ok := c.server.BufferNotiff[c.userID]; ok && len(notifications.List) > 0 {
			notifications.Lock.Lock()
			packet := notifications.List

			err := c.writeToStream(packet)
			if err != nil {
				return
			}
			notifications.List = nil
			c.server.BufferNotiff[c.userID] = notifications
			notifications.Lock.Unlock()
		}
		time.Sleep(1 * time.Second)
	}
}

//NewChannel init
func NewChannel(conn net.Conn, server *server) *Channel {

	ch := &Channel{conn: conn, server: server}
	fmt.Println(ch)
	//go c.reader()
	//go c.writer()
	return ch

}
