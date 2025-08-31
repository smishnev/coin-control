import React, { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useParams, useNavigate } from "react-router-dom";
import { GetCoinIconURLs } from "../../wailsjs/go/bybit/BybitService";

const CoinDetail: React.FC = () => {
  const { coinId } = useParams<{ coinId: string }>();
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [iconUrl, setIconUrl] = useState<string>("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const handleBack = () => {
    navigate('/bybit');
  };

  useEffect(() => {
    if (!coinId) {
      setError('Coin ID not provided');
      setLoading(false);
      return;
    }

    (async () => {
      try {
        const icons = await GetCoinIconURLs([coinId.toUpperCase()]);
        
        if (icons && icons.length > 0) {
          const isDark = document.documentElement.classList.contains('dark');
          const icon = icons[0];
          const chosen = isDark
            ? (icon.darkUrl ?? icon.iconUrl)
            : (icon.lightUrl ?? icon.iconUrl);
          setIconUrl(chosen ?? '');
        }
      } catch (e: any) {
        setError(String(e));
      } finally {
        setLoading(false);
      }
    })();
  }, [coinId]);

  if (loading) return <div>{t('Loading...')}</div>;
  if (error) return <div className="text-red-500">{error}</div>;
  if (!coinId) return <div className="text-red-500">Coin ID not found</div>;

  return (
    <div className="p-4 space-y-6">
      {/* Breadcrumbs */}
      <nav className="flex items-center space-x-2 text-sm text-gray-600 dark:text-gray-400">
        <span>Bybit</span>
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
        </svg>
        <span className="font-medium text-gray-900 dark:text-gray-100">{coinId.toUpperCase()}</span>
      </nav>

      {/* Back button */}
      <button 
        onClick={handleBack}
        className="flex items-center text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300 transition-colors"
      >
        <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
        </svg>
        {t('backToCoins')}
      </button>

      {/* Coin header */}
      <div className="flex items-center space-x-4">
        {iconUrl && (
          <img 
            src={iconUrl} 
            alt={coinId} 
            className="w-16 h-16"
          />
        )}
        <div>
          <h1 className="text-3xl font-bold">{coinId.toUpperCase()}</h1>
          <p className="text-gray-600 dark:text-gray-400">Cryptocurrency</p>
        </div>
      </div>

      {/* Price section placeholder */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <h2 className="text-xl font-semibold mb-4">{t('realTimePrice')}</h2>
        <div className="text-2xl font-bold text-green-600">
          {t('comingSoon')}
        </div>
        <p className="text-sm text-gray-500 mt-2">
          {t('priceDescription')}
        </p>
      </div>
    </div>
  );
};

export default CoinDetail;
