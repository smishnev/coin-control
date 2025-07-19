import React, { useEffect, useState } from "react";
import { GetUser, UpdateUser } from "../../wailsjs/go/user/UserService";
import { GetBybitByUserId, UpsertBybit } from "../../wailsjs/go/bybit/BybitService";
import { user } from "../../wailsjs/go/models";
import { useTranslation } from 'react-i18next';
import { useAuth } from '../contexts/AuthContext';

const UserProfileForm: React.FC = () => {
  const { t } = useTranslation();
  const { user: authUser } = useAuth();

  const [form, setForm] = useState<user.User & {bybitApiKey: string}>({
    id: "",
    firstName: "",
    lastName: "",
    bybitApiKey: "",
  });
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState("");

  useEffect(() => {
    if (authUser) { 
      GetUser(authUser.user_id)
        .then((userData) => {
          return GetBybitByUserId(authUser.user_id)
            .then((bybitData) => {
              return { userData, bybitData };
            })
            .catch((bybitError) => {
              console.log('Bybit data error (this is normal for new users):', bybitError);
              return { userData, bybitData: null };
            });
        })
        .then(({ userData, bybitData }) => {
          setForm({
            ...userData,
            bybitApiKey: bybitData && bybitData.apiKey ? bybitData.apiKey : ''
          });
          setLoading(false);
        })
        .catch((error) => {
          console.error('Error loading user data:', error);
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
    const { bybitApiKey, ...userData } = form;
    
    try {
      await UpdateUser(userData);
      await UpsertBybit(bybitApiKey, userData.id);
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
        <h2 className="text-2xl font-bold text-foreground">User Profile</h2>
        <p className="text-foreground">Welcome, {authUser.nickname}!</p>
      </div>
      
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-foreground mb-1">
            {t('firstName')}
          </label>
          <input
            type="text"
            name="firstName"
            value={form.firstName}
            onChange={handleChange}
            className="w-full text-muted-foreground bg-menu px-3 py-2 border border-border rounded-md focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
            placeholder="Enter your first name"
          />
        </div>
        
        <div>
          <label className="block text-sm font-medium text-foreground mb-1">
            {t('lastName')}
          </label>
          <input
            type="text"
            name="lastName"
            value={form.lastName}
            onChange={handleChange}
            className="w-full bg-menu text-muted-foreground px-3 py-2 border border-border rounded-md focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
            placeholder="Enter your last name"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-foreground mb-1">
            {t('bybitApiKey')}
          </label>
          <input
            type="text"
            name="bybitApiKey"
            value={form.bybitApiKey}
            onChange={handleChange}
            className="w-full bg-menu text-muted-foreground px-3 py-2 border border-border rounded-md focus:outline-none focus:ring-2 focus:ring-brand focus:border-brand"
            placeholder="Enter your Bybit key"
          />
        </div>
        
        <button 
          type="submit" 
          className="w-full bg-brand text-white px-4 py-2 rounded-md hover:brightness-90 focus:outline-none focus:ring-2 focus:ring-brand focus:ring-offset-2 transition duration-200"
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