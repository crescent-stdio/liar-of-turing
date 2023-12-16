import { Message, User } from "./playerTypes";

export interface WsJsonResponse {
  timestamp: number;
  max_player: number;
  action: string;
  user: User;
  message: string;
  message_type: string;
  message_log_list: Message[];
  online_user_list: User[];
  player_list: User[];
  game_turns_left: number;
  game_round: number;
  game_turn_num: number;
  game_round_num: number;
  is_game_started: boolean;
  is_game_over: boolean;
  user_selections_list: UserSelection[];
}
export interface WsJsonRequest {
  action: string;
  // room_id: number;
  max_player: number;
  user: User;
  timestamp: number;
  message: string;
  game_turns_left: number;
  game_round: number;
  game_turn_num: number;
  game_round_num: number;
  user_selection: UserSelection;
}

export interface UserSelection {
  timestamp: number;
  user: User;
  selection: string;
  reason: string;
}
