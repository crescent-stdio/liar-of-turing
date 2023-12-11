import { Message, Room, User } from "@/types/playerTypes";
import { atomWithReset } from "jotai/utils";

export const roomListAtom = atomWithReset<Room[]>(
  Array(1)
    .fill([])
    .map((_, idx) => {
      return {
        room_id: idx + 1,
        room_status: "waiting",
        room_members: [],
        room_messages: [],
        game_status: {
          players: [],
          isOver: false,
        },
      };
    })
);

export const currentRoomAtom = atomWithReset<Room>({
  room_id: 0,
  room_status: "waiting",
  room_members: [],
  room_messages: [],
  game_status: {
    players: [],
    isOver: false,
  },
});

export const roomInfoAtom = atomWithReset<Room[]>([
  {
    room_id: 0,
    room_status: "waiting",
    room_members: [],
    room_messages: [],
    game_status: {
      players: [],
      isOver: false,
    },
  },
]);

export const userAtom = atomWithReset<User>({
  uuid: "",
  user_id: 0,
  nickname_id: 0,
  username: "",
  role: "",
  is_online: false,
  player_type: "",
});

export const userListAtom = atomWithReset<User[]>([]);

export const messageAtom = atomWithReset<string>("");

export const messageLogAtom = atomWithReset<Message>({
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
});
export const messageLogListAtom = atomWithReset<Message[]>([]);

export const socketAtom = atomWithReset<WebSocket | null>(null);

export const wsConnectedAtom = atomWithReset<boolean>(false);
