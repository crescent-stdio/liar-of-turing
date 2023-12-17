import { userAtom } from "@/store/chatAtom";
import { initialUserSelection } from "@/store/chatStore";
import { User } from "@/types/playerTypes";
import { WsJsonRequest, WsJsonResponse } from "@/types/wsTypes";
import { init } from "next/dist/compiled/webpack/webpack";

export const sendEnterHumanByUUID = (
  socket: WebSocket | null,
  userUUID: string,
  maxPlayer: number
) => {
  const jsonData: WsJsonRequest = {
    action: "enter_human",
    max_player: maxPlayer,
    // room_id: 0,
    user: {
      uuid: userUUID,
      user_id: -1,
      nickname_id: -1,
      username: "",
      role: "human",
      is_online: true,
      player_type: "player",
    },
    timestamp: Date.now(),
    message: "",
    game_round: 0,
    game_turns_left: 0,
    game_round_num: 0,
    game_turn_num: 0,
    user_selection: initialUserSelection,
  };
  socket?.send(JSON.stringify(jsonData));
};

export const sendEnterHumanByUserData = (
  socket: WebSocket | null,
  userData: User,
  maxPlayer: number
) => {
  const jsonData: WsJsonRequest = {
    action: "enter_human",
    max_player: maxPlayer,
    user: userData,
    timestamp: Date.now(),
    message: "",
    game_round: 0,
    game_turns_left: 0,
    game_round_num: 0,
    game_turn_num: 0,
    user_selection: initialUserSelection,
  };
  socket?.send(JSON.stringify(jsonData));
};

export const sendLeftUser = (socket: WebSocket | null) => {
  const jsonData = { action: "left_user" };
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(jsonData));
  }
};
