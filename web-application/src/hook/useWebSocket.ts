import {
  chatAtom,
  chatLogAtom,
  messageLogListAtom,
  socketAtom,
  updateChatLog,
  userAtom,
  userListAtom,
} from "@/store/chatAtom";
import { initialMessage } from "@/store/chatStore";
import { maxPlayerAtom } from "@/store/gameAtom";
import { Message, User } from "@/types/playerTypes";
import { WsJsonRequest, WsJsonResponse } from "@/types/wsTypes";
import { getUserUUID } from "@/utils/liarHelper";
import {
  sendEnterHumanByUUID,
  sendEnterHumanByUserData,
  sendLeftUser,
} from "@/utils/wsHelper";
import { atom, useAtom, useAtomValue } from "jotai";
import { useResetAtom } from "jotai/utils";
import { useState, useEffect, useCallback } from "react";

const WEBSOCKET_URL = process.env.NEXT_PUBLIC_WEBSOCKET_URL || "";

export default function useWebSocket(
  userUUID: string | null,
  userData: User | null
) {
  const [socket, setSocket] = useAtom(socketAtom);
  const [isConnected, setIsConnected] = useState<boolean>(false);
  const [userList, setUserList] = useAtom(userListAtom);
  const [user, setUser] = useAtom(userAtom);
  const [messageLogList, setMessageLogList] = useAtom(messageLogListAtom);
  const [, setChatLog] = useAtom(chatLogAtom);
  const [maxPlayer, setMaxPlayer] = useAtom(maxPlayerAtom);
  // const [];

  // Function to handle incoming WebSocket messages
  const handleWebSocketMessage = useCallback((event: MessageEvent) => {
    const data: WsJsonResponse = JSON.parse(event.data);
    console.log("Received action:", data.action);

    switch (data.action) {
      case "human_info":
        setUserList(() => {
          if (!data.online_user_list) return [];
          return data.online_user_list.filter((user) => user.role !== "admin");
        });
        if (userUUID && userUUID === data.user.uuid) {
          setUser(data.user);
        }

        updateChatLog(setChatLog, data);

        setMaxPlayer(data.max_player);
        setMessageLogList(() => {
          if (!data.message_log_list) return [];
          return data.message_log_list;
        });

        break;
      case "user_list":
        console.log("user_list", data.online_user_list);

        setMaxPlayer(data.max_player);
        setUserList(() => {
          if (!data.online_user_list) return [];
          return data.online_user_list.filter((user) => user.role !== "admin");
        });

        break;
      case "new_message":
        console.log("Message", data);
        updateChatLog(setChatLog, data);

        setMaxPlayer(data.max_player);
        setMessageLogList(() => {
          if (!data.message_log_list) return [];
          return data.message_log_list;
        });

        break;
      case "update_state":
        console.log("update_state", data);
        if (data.user.uuid === userUUID) {
          setUser(data.user);
        }
        setMaxPlayer(data.max_player);
        setUserList(() => {
          if (!data.online_user_list) return [];
          return data.online_user_list.filter((user) => user.role !== "admin");
        });
        updateChatLog(setChatLog, data);
        setMessageLogList(() => {
          if (!data.message_log_list) return [];
          return data.message_log_list;
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
        if (userUUID) sendEnterHumanByUUID(socket, userUUID, maxPlayer);
        else if (userData)
          sendEnterHumanByUserData(socket, userData, maxPlayer);
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
