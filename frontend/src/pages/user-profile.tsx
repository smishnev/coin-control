import React, { useEffect, useState } from "react";
import { GetUser, CreateOrUpdate } from "../../wailsjs/go/user/UserService";
import { user } from "../../wailsjs/go/models";
import { useTranslation } from 'react-i18next';

const UserProfileForm: React.FC = () => {
  const { t } = useTranslation();

  const [form, setForm] = useState<user.User>({
    id: "",
    firstName: "",
    lastName: "",
  });
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState("");

  useEffect(() => {
    //TODO Use real user ID from your authentication context or state
    const userId = "58641274-0ce0-4e59-ab11-7cbfcb4e9370";
    GetUser(userId)
      .then((data) => {
        setForm(data);
        setLoading(false);
      })
      .catch(() => setLoading(false));
  }, []);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await CreateOrUpdate(form);
      setMessage(t('userProfileUpdateSuccess'));
    } catch {
      setMessage(t('userProfileUpdateError'));
    }
  };

  if (loading) return <div>{t('loading')}</div>;

  return (
    <form onSubmit={handleSubmit} className="max-w-md mx-auto space-y-4">
      <div>
        <label>{t('firstName')}</label>
        <input
          type="text"
          name="firstName"
          value={form.firstName}
          onChange={handleChange}
          className="border p-2 w-full"
        />
      </div>
      <div>
        <label>{t('lastName')}</label>
        <input
          type="text"
          name="lastName"
          value={form.lastName}
          onChange={handleChange}
          className="border p-2 w-full"
        />
      </div>
      <button type="submit" className="bg-blue-500 text-white px-4 py-2 rounded">
        {t('saveButton')}
      </button>
      {message && <div className="mt-2">{message}</div>}
    </form>
  );
};

export default UserProfileForm;