import { WsJsonRequest } from "@/types/wsTypes";
import MessageInput from "../MessageInput";
import HorizontalLine from "../Line/HorizontalLine";
import { useAtom } from "jotai";
import { userAtom } from "@/store/chatAtom";

type InputModalProps = {
  isGameStarted: boolean;
  isYourTurn: boolean;
  sendMessage: (message: WsJsonRequest) => void;
};
export default function InputModal({
  isGameStarted,
  isYourTurn,
  sendMessage: handleWebSocketMessageSend,
}: InputModalProps) {
  const [user, setUser] = useAtom(userAtom);

  if (!user) return <></>;
  if (isGameStarted && user.player_type === "player" && isYourTurn) {
    return (
      <MessageInput
        userData={user}
        key={user.uuid}
        sendMessage={handleWebSocketMessageSend}
      />
    );
  }
  return <HorizontalLine />;
}
