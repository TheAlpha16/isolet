import { TeamType } from "@/utils/types";

interface ProfileStore {
	team: TeamType;
	fetchProfile: () => void;
}
