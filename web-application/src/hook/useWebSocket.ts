import {
  chatAtom,
  chatLogAtom,
  messageLogListAtom,
  playerListAtom,
  socketAtom,
  updateChatLog,
  userAtom,
  userListAtom,
} from "@/store/chatAtom";
import { initialMessage } from "@/store/chatStore";
import {
  gameRoundAtom,
  gameTurnsLeftAtom,
  isFinishedRoundAtom,
  isFinishedSubmitionAtom,
  isGameStartedAtom,
  isYourTurnAtom,
  maxPlayerAtom,
} from "@/store/gameAtom";
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
  const [, setIsYourTurn] = useAtom(isYourTurnAtom);
  // const [];
  const [gameTurnsLeft, setGameTurnsLeft] = useAtom(gameTurnsLeftAtom);
  const [gameRound, setGameRound] = useAtom(gameRoundAtom);
  const [playerList, setPlayerList] = useAtom(playerListAtom);
  const [isGameStarted, setIsGameStarted] = useAtom(isGameStartedAtom);
  const [isFinishedRound, setIsFinishedRound] = useAtom(isFinishedRoundAtom);
  const [isFinishedSubmition, setIsFinishedSubmition] = useAtom(
    isFinishedSubmitionAtom
  );
  const [, setGameRoundNum] = useAtom(gameRoundAtom);
  const [, setGameTurnsNum] = useAtom(gameTurnsLeftAtom);

  // Function to handle incoming WebSocket messages
  const handleWebSocketMessage = useCallback((event: MessageEvent) => {
    const data: WsJsonResponse = JSON.parse(event.data);
    console.log("Received action:", data.action);

    switch (data.action) {
      case "choose_ai":
        console.log("choose_ai", data);

        setIsGameStarted(false);
        setIsFinishedSubmition(false);
        setIsFinishedRound(true);

        setUserList(() => {
          if (!data.online_user_list) return [];
          return data.online_user_list.filter((user) => user.role !== "admin");
        });
        setMessageLogList(() => {
          if (!data.message_log_list) return [];
          return data.message_log_list;
        });
        break;
      case "wait_for_players":
        setIsFinishedSubmition(true);
      case "submit_ai":
        setIsFinishedSubmition(false);
      case "human_info":
        setUserList(() => {
          if (!data.online_user_list) return [];
          return data.online_user_list.filter((user) => user.role !== "admin");
        });
        setMessageLogList(() => {
          if (!data.message_log_list) return [];
          return data.message_log_list;
        });
        setPlayerList(() => {
          if (!data.player_list) return [];
          return data.player_list;
        });
        setGameRound(data.game_round);
        setGameTurnsLeft(data.game_turns_left);
        setMessageLogList(() => {
          if (!data.message_log_list) return [];
          return data.message_log_list;
        });

        break;
      case "user_list":
        console.log("user_list", data.online_user_list);
        setUserList(() => {
          if (!data.online_user_list) return [];
          return data.online_user_list.filter((user) => user.role !== "admin");
        });

        break;
      case "new_message":
        console.log("Message", data);
        setMessageLogList(() => {
          if (!data.message_log_list) return [];
          return data.message_log_list;
        });

        break;
      case "your_turn":
        console.log("your_turn", data);
        setIsYourTurn(true);
        setIsGameStarted(true);
        setMessageLogList((prevMessageLogList) => {
          if (!data.message_log_list) return [];
          const messageLog: Message = {
            timestamp: data.timestamp,
            message_id: data.message_id,
            user: data.user,
            message: data.message,
            message_type: data.message_type,
          };
          return [...prevMessageLogList, messageLog];
        });
        break;
      case "update_state":
        console.log("update_state", data);
        setUserList(() => {
          if (!data.online_user_list) return [];
          return data.online_user_list.filter((user) => user.role !== "admin");
        });
        setMessageLogList(() => {
          if (!data.message_log_list) return [];
          return data.message_log_list;
        });
        setGameRoundNum(data.game_round);
        setGameTurnsNum(data.game_turns_left);
        break;
      case "round_start":
        setIsGameStarted(true);
        setIsFinishedRound(false);
        break;
      case "next_round":
        console.log("next_round", data);
        setIsFinishedRound(false);
        setIsFinishedSubmition(false);
        setIsYourTurn(false);
        setMessageLogList(() => {
          if (!data.message_log_list) return [];
          return data.message_log_list;
        });
        setPlayerList(data.player_list);
        break;
      case "restart_game":
        console.log("restart_game", data);
        setIsFinishedRound(false);
        setIsFinishedSubmition(false);
        setIsYourTurn(false);
        setMessageLogList(() => {
          if (!data.message_log_list) return [];
          return data.message_log_list;
        });
        setGameRound(data.game_round);
        setGameTurnsLeft(data.game_turns_left);
        setPlayerList(data.player_list);
        break;
      case "restart_round":
        console.log("restart_round", data);
        setIsFinishedRound(false);
        setIsFinishedSubmition(false);
        setIsYourTurn(false);
        setMessageLogList(() => {
          if (!data.message_log_list) return [];
          return data.message_log_list;
        });
        setPlayerList(data.player_list);
        setGameRound(data.game_round);
        setGameTurnsLeft(data.game_turns_left);
        break;
      case "game_over":
        console.log("game_over", data);
        setIsFinishedRound(false);
        setIsGameStarted(false);
        setIsYourTurn(false);
        setMessageLogList(() => {
          if (!data.message_log_list) return [];
          return data.message_log_list;
        });
        break;
    }

    if (userUUID) {
      const myUser = data.online_user_list.find(
        (user: User) => user.uuid === userUUID
      );
      if (myUser) {
        setUser(myUser);
      }
    }
    console.log("data", data);
    updateChatLog(setChatLog, data);
    if (data.game_turns_left >= 0) setGameTurnsLeft(data.game_turns_left);
    if (data.game_round > 0) setGameRound(data.game_round);
    if (data.max_player) setMaxPlayer(data.max_player);
    if (data.player_list && data.player_list.length >= 0)
      setPlayerList(data.player_list);
    // }, []);
  }, []); // Add other dependencies as needed

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
  }, [
    socket,
    messageLogList,
    userUUID,
    playerList,
    userData,
    maxPlayer,
    gameTurnsLeft,
    gameRound,
    isFinishedRound,
    isFinishedSubmition,
  ]);

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
