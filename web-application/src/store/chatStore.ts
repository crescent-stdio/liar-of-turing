import { Message, User } from "@/types/playerTypes";

// export const
export const initialMessage: Message = {
  timestamp: 0,
  message_id: 0,
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
