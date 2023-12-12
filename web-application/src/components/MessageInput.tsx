import useMessageInput from "@/hook/useMessageInput";
import {
  gameRoundAtom,
  gameTurnsLeftAtom,
  isYourTurnAtom,
  maxPlayerAtom,
} from "@/store/gameAtom";
import { User } from "@/types/playerTypes";
import { WsJsonRequest } from "@/types/wsTypes";
import { useAtom, useAtomValue } from "jotai";

type MessageInputProps = {
  userData: User;
  sendMessage: (message: WsJsonRequest) => void;
};

const MessageInput: React.FC<MessageInputProps> = ({
  userData,
  sendMessage,
}) => {
  const {
    // inputMessage,
    message,
    handleMessageChange,
    handleSubmit: handleCustomSubmit,
    resetMessage,
  } = useMessageInput();
  const maxPlayer = useAtomValue(maxPlayerAtom);
  const [, setIsYourTurn] = useAtom(isYourTurnAtom);
  const gameTurnsLeft = useAtomValue(gameTurnsLeftAtom);
  const gameRound = useAtomValue(gameRoundAtom);

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    handleCustomSubmit(event);
    const jsonData: WsJsonRequest = {
      max_player: maxPlayer,
      action: "new_message",
      user: userData,
      timestamp: Date.now(),
      message: message,
      game_round: gameRound,
      game_turns_left: gameTurnsLeft,
    };
    sendMessage(jsonData);
    resetMessage();
    setIsYourTurn(false);
  };

  return (
    <form className="mt-4 flex flex-row" onSubmit={handleSubmit}>
      <label htmlFor="message">
        {userData.username && (
          <span
            className="mr-2 font-bold flex-1"
            style={{
              color: "#3b82f6",
            }}
          >{`${userData.username}: `}</span>
        )}
      </label>
      <input
        autoFocus
        className="border-2 border-gray-400 rounded-md flex-1"
        type="text"
        id="message"
        value={message}
        onChange={handleMessageChange}
      />
    </form>
  );
};

export default MessageInput;
