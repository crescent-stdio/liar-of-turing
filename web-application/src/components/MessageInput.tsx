import useMessageInput from "@/hook/useMessageInput";
import { User } from "@/types/playerTypes";
import { WsJsonRequest } from "@/types/wsTypes";

type MessageInputProps = {
  userData: User;
  onSendMessage: (message: WsJsonRequest) => void;
};

const MessageInput: React.FC<MessageInputProps> = ({
  userData,
  onSendMessage,
}) => {
  const {
    // inputMessage,
    message,
    handleMessageChange,
    handleSubmit: handleCustomSubmit,
    resetMessage,
  } = useMessageInput();

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    handleCustomSubmit(event);
    const jsonData: WsJsonRequest = {
      action: "new_message",
      user: userData,
      timestamp: Date.now(),
      message: message,
    };
    onSendMessage(jsonData);
    resetMessage();
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