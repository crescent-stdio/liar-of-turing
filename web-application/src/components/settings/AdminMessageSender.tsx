import { useState } from "react";
import { User } from "@/types/playerTypes";
import { WsJsonRequest } from "@/types/wsTypes";
import { useAtomValue } from "jotai";
import { maxPlayerAtom } from "@/store/gameAtom";
import { initialUserSelection } from "@/store/chatStore";

type AdminMessageSenderProps = {
  userList: User[];
  sendMessage: (message: WsJsonRequest) => void;
};

const AdminMessageSender: React.FC<AdminMessageSenderProps> = ({
  userList,
  sendMessage,
}) => {
  const [isDebugMode, setIsDebugMode] = useState<boolean>(false);
  const [testUsername, setTestUsername] = useState<string>(
    userList[0].username
  );
  const [testMessage, setTestMessage] = useState<string>("");
  const maxPlayer = useAtomValue(maxPlayerAtom);

  const handleTestSendMessage = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const username = testUsername;
    const message = testMessage;
    if (username.length > 0 && message.length > 0) {
      const testUserData = userList.find(
        (user: User) => user.username === username
      );
      if (!testUserData) return;
      const jsonData: WsJsonRequest = {
        max_player: maxPlayer,
        action: "new_message_admin",
        user: testUserData,
        timestamp: Date.now(),
        message: testMessage,
        game_round: 0,
        game_turns_left: 0,
        game_round_num: 0,
        game_turn_num: 0,
        user_selection: initialUserSelection,
      };
      sendMessage(jsonData);
      setTestMessage("");
    }
  };

  return (
    <div className="mt-4">
      {/* <button
        className="top-0 left-0 absolute text-white hover:text-black"
        onClick={() => setIsDebugMode((isDebugMode) => !isDebugMode)}
      >
        {isDebugMode ? "Hide" : "Show"} sending message form
      </button>
      {isDebugMode && ( */}
      <div className="flex flex-col">
        <h3 className="mt-4 font-bold text-xl">Send message</h3>
        <form className="flex flex-row" onSubmit={handleTestSendMessage}>
          <label htmlFor="username">Username</label>
          <select
            name="username"
            id="username"
            onChange={(e) => {
              setTestUsername(e.target.value);
            }}
            className="border-2 border-gray-400 rounded-md w-fit-content"
          >
            {userList &&
              userList.length > 0 &&
              userList.map((user: User, index) => {
                return (
                  <option key={index} value={user.username}>
                    {user.username}
                  </option>
                );
              })}
          </select>
          <label htmlFor="message" className="mx-2 ">
            Message
          </label>
          <input
            className="border-2 border-gray-400 rounded-md w-fit-content"
            type="text"
            id="message"
            value={testMessage}
            onChange={(e) => setTestMessage(e.target.value)}
          />
        </form>
      </div>
    </div>
  );
};

export default AdminMessageSender;
