import React, { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";
import { FetchSpotHoldings, GetCoinIconURLs, PrefetchCoinIcons } from "../../wailsjs/go/main/App";
import { useAuth } from "../contexts/AuthContext";

type Holding = { coin: string; free: string; locked: string };

const BybitForm: React.FC = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { user: authUser } = useAuth();
  const [holdings, setHoldings] = useState<Holding[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [iconUrls, setIconUrls] = useState<Record<string,string>>({});

  useEffect(() => {
    if (!authUser) {
      setLoading(false);
      return;
    }

    (async () => {
      try {
        const data = await FetchSpotHoldings(authUser.user_id);
        setHoldings(data);

        // Fetch icons for coins
        const coins = Array.from(new Set(data.map(h => h.coin.toUpperCase())));
        const icons = await GetCoinIconURLs(coins);

        const iconMap: Record<string,string> = {};
        const isDark = document.documentElement.classList.contains('dark');
        (icons || []).forEach((i: any) => {
          const chosen = isDark
            ? (i.darkUrl ?? i.iconUrl)
            : (i.lightUrl ?? i.iconUrl);
          iconMap[i.coin.toUpperCase()] = chosen ?? '';
        });
        setIconUrls(iconMap);

        // Warm local disk cache in background (no await)
        PrefetchCoinIcons(coins).catch(() => {});
      } catch (e: any) {
        setError(String(e));
      } finally {
        setLoading(false);
      }
    })();
  }, [authUser]);

  // Filter only coins with positive balance
  const coinsWithBalance = holdings.filter(h => {
    const qty = (parseFloat(h.free) || 0) + (parseFloat(h.locked) || 0);
    return qty > 0;
  });

  if (loading) return <div>{t('Loading...')}</div>;
  if (error) return <div className="text-red-500">{error}</div>;

  return (
    <div className="p-4 space-y-4">
      <h1 className="text-2xl font-bold">{t('myCoins')}</h1>
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
        {coinsWithBalance.map((holding) => {
          const icon = iconUrls[holding.coin.toUpperCase()];
          return (
            <div 
              key={holding.coin} 
              className="flex flex-col items-center p-4 border border-border rounded-lg hover:bg-gray-50 dark:hover:bg-gray-800 cursor-pointer transition-colors"
              onClick={() => navigate(`/bybit/${holding.coin.toLowerCase()}`)}
            >
              {icon && <img src={icon} alt={holding.coin} className="w-8 h-8 mb-2" />}
              <span className="text-sm font-medium">{holding.coin.toUpperCase()}</span>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default BybitForm;