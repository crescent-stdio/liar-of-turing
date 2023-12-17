import { Message, User } from "@/types/playerTypes";
import { UserSelection, WsJsonRequest } from "@/types/wsTypes";

// export const
export const initialMessage: Message = {
  timestamp: 0,
  user: {
    uuid: "",
    user_id: 0,
    nickname_id: 0,
    username: "",
    role: "",
    is_online: false,
    player_type: "",
  },
  message: "",
  message_type: "",
};

export const adminUser: User = {
  uuid: "0",
  user_id: 0,
  nickname_id: 999,
  username: "server",
  role: "admin",
  is_online: false,
  player_type: "admin",
};

export const initialUserSelection: UserSelection = {
  timestamp: 0,
  user: {
    uuid: "",
    user_id: 0,
    nickname_id: 0,
    username: "",
    role: "",
    is_online: false,
    player_type: "",
  },
  selection: "",
  reason: "",
};

export const initialWsJsonRequest: WsJsonRequest = {
  max_player: 0,
  action: "",
  user: {
    uuid: "",
    user_id: 0,
    nickname_id: 0,
    username: "",
    role: "",
    is_online: false,
    player_type: "",
  },
  timestamp: 0,
  message: "",
  game_turns_left: 0,
  game_round: 0,
  game_turn_num: 0,
  game_round_num: 0,
  user_selection: initialUserSelection,
};
