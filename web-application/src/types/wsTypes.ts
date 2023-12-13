import { Message, User } from "./playerTypes";

export interface WsJsonResponse {
  timestamp: number;
  max_player: number;
  message_id: number;
  action: string;
  user: User;
  message: string;
  message_type: string;
  message_log_list: Message[];
  online_user_list: User[];
  player_list: User[];
  game_round: number;
  game_turns_left: number;
  is_game_started: boolean;
  is_game_over: boolean;
  user_selections_list: UserSelection[];
}
export interface WsJsonRequest {
  action: string;
  max_player: number;
  // room_id: number;
  user: User;
  timestamp: number;
  message: string;
  game_round: number;
  game_turns_left: number;
  user_selection: UserSelection;
}

export interface UserSelection {
  user: User;
  selection: string;
  reason: string;
}
