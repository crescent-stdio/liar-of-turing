// import { User } from "./wsTypes";

export interface User {
  uuid: string;
  user_id: number;
  nickname_id: number;
  username: string;
  role: string;
  is_online: boolean;
  player_type: string;
}

export interface Game {
  players: User[];
  isOver: boolean;
}
export interface Message {
  timestamp: number;
  message_id: number;
  user: User;
  message: string;
  message_type: string;
}

export interface Room {
  room_id: number;
  // room_name: string;
  // room_type: string;
  room_status: string;
  room_members: User[];
  room_messages: Message[];
  game_status: Game;
}
