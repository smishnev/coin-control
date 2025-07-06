import React, { useEffect, useState } from "react";
import { GetUser, CreateOrUpdate } from "../../wailsjs/go/user/UserService";
import { user } from "../../wailsjs/go/models";
import { useTranslation } from 'react-i18next';
import { useAuth } from '../contexts/AuthContext';

const UserProfileForm: React.FC = () => {
  const { t } = useTranslation();
  const { user: authUser } = useAuth();

  const [form, setForm] = useState<user.User>({
    id: "",
    firstName: "",
    lastName: "",
  });
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState("");

  useEffect(() => {
    if (authUser) {
      GetUser(authUser.user_id)
        .then((data) => {
          setForm(data);
          setLoading(false);
        })
        .catch((error) => {
          console.error('Error loading user profile:', error);
          setLoading(false);
        });
    } else {
      setLoading(false);
    }
  }, [authUser]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await CreateOrUpdate(form);
      setMessage(t('userProfileUpdateSuccess'));
    } catch (error) {
      console.error('Error updating user profile:', error);
      setMessage(t('userProfileUpdateError'));
    }
  };

  if (loading) return <div>{t('loading')}</div>;

  if (!authUser) {
    return <div className="text-center text-red-600">User not authenticated</div>;
  }

  return (
    <div className="max-w-md mx-auto space-y-6">
      <div className="text-center">
        <h2 className="text-2xl font-bold text-gray-900">User Profile</h2>
        <p className="text-gray-600">Welcome, {authUser.nickname}!</p>
      </div>
      
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            {t('firstName')}
          </label>
          <input
            type="text"
            name="firstName"
            value={form.firstName}
            onChange={handleChange}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            placeholder="Enter your first name"
          />
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            {t('lastName')}
          </label>
          <input
            type="text"
            name="lastName"
            value={form.lastName}
            onChange={handleChange}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            placeholder="Enter your last name"
          />
        </div>
        
        <button 
          type="submit" 
          className="w-full bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-colors"
        >
          {t('saveButton')}
        </button>
        
        {message && (
          <div className={`mt-2 p-3 rounded-md text-sm ${
            message.includes('success') 
              ? 'bg-green-100 text-green-700 border border-green-200' 
              : 'bg-red-100 text-red-700 border border-red-200'
          }`}>
            {message}
          </div>
        )}
      </form>
    </div>
  );
};

export default UserProfileForm;