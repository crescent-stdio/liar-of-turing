import { Message, Room, User } from "@/types/playerTypes";
import { SetStateAction, WritableAtom, atom, useAtom, useSetAtom } from "jotai";
import { atomWithReset } from "jotai/utils";
import { initialMessage } from "./chatStore";
import { Update } from "next/dist/build/swc";
import { WsJsonResponse } from "@/types/wsTypes";

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

export const userAtom = atom<User>({
  uuid: "",
  user_id: 0,
  nickname_id: 0,
  username: "",
  role: "",
  is_online: false,
  player_type: "",
});

export const userListAtom = atom<User[]>([]);
export const playerListAtom = atom<User[]>([]);

export const watcherListAtom = atom<User[]>((get: any) => {
  const userList = get(userListAtom);
  const playerList = get(playerListAtom);
  const watcherList = userList.filter(
    (user: User) =>
      !playerList.some((player: User) => player.uuid === user.uuid)
  );

  return watcherList;
});

export const chatAtom = atomWithReset<string>("");

export const chatLogAtom = atom<Message>(initialMessage);
export const updateChatLog: (set: any, data: WsJsonResponse | null) => void = (
  set,
  data
) => {
  if (!data) {
    set(chatLogAtom, initialMessage);
  } else {
    const message = {
      timestamp: data.timestamp,
      user: data.user,
      message: data.message,
      message_type: data.message_type,
    };
    set(chatLogAtom, message);
  }
};

export const chatLogListAtom = atomWithReset<Message[]>([]);

export const messageLogListAtom = atomWithReset<Message[]>([]);

export const socketAtom = atomWithReset<WebSocket | null>(null);

export const wsConnectedAtom = atomWithReset<boolean>(false);
