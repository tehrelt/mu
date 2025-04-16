import { Profile } from "~/schemes/profile";

type UserCardProps = {
  profile: Profile;
};
export const UserCard = ({ profile }: UserCardProps) => {
  return (
    <div className="user-card">
      <h2>
        {profile.lastName} {profile.firstName}
      </h2>
      <p>{profile.email}</p>
    </div>
  );
};
