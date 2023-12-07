export interface WsJsonResponse {
  timestamp: number;
  message_id: number;
  action: string;
  user: User;
  message: string;
  message_type: string;
  online_user_list: User[];
}
export interface WsJsonRequest {
  action: string;
  room_id: number;
  user: User;
  timestamp: number;
  message: string;
}

export interface User {
  uuid: string;
  user_id: number;
  room_id: number;
  nickname_id: number; // TODO: It will be deprecated?
  username: string;
  role: string;
  is_online: boolean;
}