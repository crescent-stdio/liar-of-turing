import { User } from "@/types/playerTypes";
import OnlineUserList from "./OnlineUserList";
import { playerListAtom, userAtom, watcherListAtom } from "@/store/chatAtom";
import { useAtomValue } from "jotai";
import { maxPlayerAtom } from "@/store/gameAtom";
import { useEffect, useState } from "react";
import HorizontalLine from "./Line/HorizontalLine";
type PlayAndWaitUserListProps = {
  // userData: User;
};
export default function PlayAndWaitUserList({}: // userData,
PlayAndWaitUserListProps) {
  const playerList = useAtomValue(playerListAtom);
  const watcherList = useAtomValue(watcherListAtom);
  const maxPlayer = useAtomValue(maxPlayerAtom);
  const [playerListTitle, setPlayerListTitle] = useState<string>("Players");
  const userData = useAtomValue(userAtom);

  // const playerListTitle = `Players [${playerList.length}/${maxPlayer}]`;
  useEffect(() => {
    setPlayerListTitle(`Players [${playerList.length}/${maxPlayer}]`);
  }, [playerList.length, maxPlayer]);
  return (
    <div className="flex md:flex-col h-full md:ml-4 md:min-h-[70vh] md:max-h-[70vh] md:w-36 flex-row max-h-full min-h-full">
      <OnlineUserList
        name={playerListTitle}
        userData={userData}
        userList={playerList}
      />
      <HorizontalLine />
      <OnlineUserList
        name={"Watchers"}
        userData={userData}
        userList={watcherList}
      />
    </div>
  );
}
