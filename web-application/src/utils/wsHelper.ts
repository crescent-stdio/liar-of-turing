import { WsJsonRequest } from "@/types/wsTypes";

export const sendEnterHuman = (socket: WebSocket | null, userUUID: string) => {
  const jsonData: WsJsonRequest = {
    action: "enter_human",
    // room_id: 0,
    user: {
      uuid: userUUID,
      user_id: -1,
      nickname_id: -1,
      username: "",
      role: "human",
      is_online: true,
      player_type: "player",
    },
    timestamp: Date.now(),
    message: "",
  };
  socket?.send(JSON.stringify(jsonData));
};

export const sendLeftUser = (socket: WebSocket | null) => {
  const jsonData = { action: "left_user" };
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(jsonData));
  }
};
