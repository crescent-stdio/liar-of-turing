import { User } from "./wsTypes";

export interface Player {
  uuid: string;
  user_id: number;
  username: string;
  role: string;
}

export interface Game {
  players: Player[];
  isOver: boolean;
}
export interface Message {
  message_id: number;
  timestamp: number;
  user: User;
  message: string;
  message_type: string;
}

export interface Room {
  room_id: number;
  room_name: string;
  room_type: string;
  room_status: string;
  room_owner: User;
  room_members: User[];
  room_messages: Message[];
  room_game: Game;
}
