import React from 'react';

const Avatar: React.FC = () => {
    return (
        <div className="flex items-center">
            <img
                src="https://via.placeholder.com/40" // Placeholder image for the avatar
                alt="User Avatar"
                className="rounded-full w-10 h-10"
            />
        </div>
    );
};

export default Avatar;