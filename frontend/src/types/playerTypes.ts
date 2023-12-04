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
  id: number;
  timestamp: number;
  sender: string;
  content: string;
}
export interface Chat {
  messages: Message[];
}
