import { User } from "@/types/playerTypes";
import OnlineUserList from "./OnlineUserList";
import { playerListAtom, watcherListAtom } from "@/store/chatAtom";
import { useAtom } from "jotai";
import { MAX_PLAYER } from "@/store/gameStore";
type PlayAndWaitUserListProps = {
  userData: User;
};
export default function PlayAndWaitUserList({
  userData,
}: PlayAndWaitUserListProps) {
  const [playerList] = useAtom(playerListAtom);
  const [watcherList] = useAtom(watcherListAtom);
  const playerListTitle = `Players [${playerList.length}/${MAX_PLAYER}]`;
  return (
    <div className="flex flex-col h-full ml-4">
      <OnlineUserList
        name={playerListTitle}
        userData={userData}
        userList={playerList}
      />
      <hr className="my-2" />
      <OnlineUserList
        name={"Watchers"}
        userData={userData}
        userList={watcherList}
      />
    </div>
  );
}
