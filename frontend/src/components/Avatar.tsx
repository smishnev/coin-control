import React, { useEffect, useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { GetUser } from "../../wailsjs/go/user/UserService";

const Avatar: React.FC = () => {
    const { user: authUser } = useAuth();
    const [initial, setInitial] = useState<string>('NA');
    const avatarUrl = ''; 

    useEffect(() => {
        const fetchInitials = async () => {
            if (!authUser) {
                setInitial('NA');
                return;
            }
            const user = await GetUser(authUser.user_id);
            const firstInitial = user.firstName?.[0]?.toUpperCase() ?? '';
            const lastInitial = user.lastName?.[0]?.toUpperCase() ?? '';
            setInitial(`${firstInitial}${lastInitial}`);
        };
        fetchInitials();
    }, [authUser]);

    return (
        <div className="flex items-center">
            {avatarUrl ? (
                <img
                    src={avatarUrl}
                    alt="User Avatar"
                    className="rounded-full w-10 h-10"
                />
            ) : (
                <div className='w-[37.5px] h-[37.5px] rounded-full bg-neutral-400 dark:bg-slate-400 flex items-center justify-center text-white font-medium'>
                    {initial}
                </div>
            )}
        </div>
    );
};

export default Avatar;