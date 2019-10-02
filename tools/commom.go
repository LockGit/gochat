/**
 * Created by lock
 * Date: 2019-08-18
 * Time: 18:03
 */
package tools

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/bwmarrin/snowflake"
	"io"
)

const SessionPrefix = "sess_"

func GetSnowflakeId() string {
	//default node id eq 1,this can modify to different serverId node
	node, _ := snowflake.NewNode(1)
	// Generate a snowflake ID.
	id := node.Generate().String()
	return id
}

func GetRandomToken(prefix string, length int) string {
	r := make([]byte, length)
	io.ReadFull(rand.Reader, r)
	return prefix + "_" + base64.URLEncoding.EncodeToString(r)
}

func CreateSessionId() string {
	return GetRandomToken(SessionPrefix, 32)
}

func GetSessionName(authToken string) string {
	return SessionPrefix + authToken
}
