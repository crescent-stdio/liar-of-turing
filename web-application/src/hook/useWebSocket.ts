import {
  messageLogListAtom,
  socketAtom,
  userAtom,
  userListAtom,
} from "@/store/chatAtom";
import { Message } from "@/types/playerTypes";
import { WsJsonRequest, WsJsonResponse } from "@/types/wsTypes";
import { getUserUUID } from "@/utils/liarHelper";
import { sendEnterHuman, sendLeftUser } from "@/utils/wsHelper";
import { useAtom } from "jotai";
import { useState, useEffect, useCallback } from "react";

const WEBSOCKET_URL = process.env.NEXT_PUBLIC_WEBSOCKET_URL || "";

export default function useWebSocket(userUUID: string) {
  const [socket, setSocket] = useAtom(socketAtom);
  const [isConnected, setIsConnected] = useState<boolean>(false);
  const [userList, setUserList] = useAtom(userListAtom);
  const [user, setUser] = useAtom(userAtom);
  const [messageLogList, setMessageLogList] = useAtom(messageLogListAtom);

  // Function to handle incoming WebSocket messages
  const handleWebSocketMessage = useCallback((event: MessageEvent) => {
    const data: WsJsonResponse = JSON.parse(event.data);
    console.log("Received action:", data.action);

    switch (data.action) {
      case "human_info":
        setUserList(() => {
          if (!data.online_user_list) return [];
          return data.online_user_list;
        });
        if (userUUID === data.user.uuid) {
          setUser(data.user);
        }
        setMessageLogList((messageLogList: Message[]) => {
          const message = {
            message_id: data.message_id,
            timestamp: data.timestamp,
            user: data.user,
            message: data.message,
            message_type: data.message_type,
          };
          return [...messageLogList, message];
        });

        break;
      case "user_list":
        console.log("user_list", data.online_user_list);

        setUserList(() => {
          if (!data.online_user_list) return [];
          return data.online_user_list;
        });

        break;
      case "new_message":
        console.log("Message", data);
        setMessageLogList((messageLogList: Message[]) => {
          const message = {
            message_id: data.message_id,
            timestamp: data.timestamp,
            user: data.user,
            message: data.message,
            message_type: data.message_type,
          };
          return [...messageLogList, message];
        });

        break;
      case "update_state":
        console.log("update_state", data);
        setUserList(() => {
          if (!data.online_user_list) return [];
          return data.online_user_list;
        });
        setMessageLogList((messageLog: Message[]) => {
          const message = {
            message_id: data.message_id,
            timestamp: data.timestamp,
            user: data.user,
            message: data.message,
            message_type: data.message_type,
          };
          return [...messageLog, message];
        });

        break;
    }
  }, []);

  // Function to send WebSocket messages
  const handleWebSocketMessageSend = useCallback(
    (message: WsJsonRequest) => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify(message));
      } else {
        console.error("WebSocket is not connected.");
      }
    },
    [socket]
  );

  const handleWebSocketClose = () => {
    console.log("WebSocket closed");
    sendLeftUser(socket);
  };

  useEffect(() => {
    if (socket) {
      socket.onopen = () => {
        console.log("connected!!");
        sendEnterHuman(socket, userUUID);
      };
      socket.onmessage = handleWebSocketMessage;
      socket.onclose = handleWebSocketClose;
    }
  }, [socket, messageLogList, userUUID]);

  useEffect(() => {
    const sc = new WebSocket(WEBSOCKET_URL);
    setSocket(sc);
    setIsConnected(true);
    return () => {};
  }, []);

  return {
    socket,
    isConnected,
    userList,
    user,
    messageLogList,
    handleWebSocketMessageSend,
  };
}
