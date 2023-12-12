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
}
export interface WsJsonRequest {
  max_player: number;
  action: string;
  // room_id: number;
  user: User;
  timestamp: number;
  message: string;
}

// export interface User {
//   uuid: string;
//   user_id: number;
//   room_id: number;
//   nickname_id: number; // TODO: It will be deprecated?
//   username: string;
//   role: string;
//   is_online: boolean;
//   status: string;
// }
