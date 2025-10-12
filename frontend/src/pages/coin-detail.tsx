import React, { useEffect, useState, useRef } from "react";
import { useTranslation } from "react-i18next";
import { useParams, useNavigate } from "react-router-dom";
import { GetCoinIconURLs, GetCurrentPrice, StartPriceStream, StopPriceStream, GetAssetBalance } from "../../wailsjs/go/main/App";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { useAuth } from "../contexts/AuthContext";

type CoinBalance = { 
  coin: string; 
  walletBalance: string; 
  transferBalance: string; 
  bonus: string; 
  locked: string
};

// Component for displaying coin balance
const CoinBalanceSection: React.FC<{ coinSymbol: string }> = ({ coinSymbol }) => {
  const { user: authUser } = useAuth();
  const [balance, setBalance] = useState<CoinBalance | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!authUser || !coinSymbol) return;

    (async () => {
      try {
        const balance = await GetAssetBalance(authUser.user_id, coinSymbol);
        setBalance(balance || null);
      } catch (e: any) {
        console.error('Failed to fetch balance:', e);
        setError(String(e));
      } finally {
        setLoading(false);
      }
    })();
  }, [authUser, coinSymbol]);

  if (!authUser) return null;
  if (loading) return <div className="text-sm text-gray-500">Loading balance...</div>;
  if (error) return <div className="text-sm text-red-500">Failed to load balance</div>;
  if (!balance) return <div className="text-sm text-gray-500">No balance data</div>;

  const walletBalance = parseFloat(balance.walletBalance) || 0;
  const transferBalance = parseFloat(balance.transferBalance) || 0;
  const locked = parseFloat(balance.locked) || 0;
  const bonus = parseFloat(balance.bonus) || 0;

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg p-4 shadow-md">
      <h3 className="text-lg font-semibold mb-3">Asset Balance</h3>
      <div className="space-y-2 text-sm">
        <div className="flex justify-between">
          <span className="text-gray-600 dark:text-gray-400">Wallet Balance:</span>
          <span className="font-mono">{walletBalance.toFixed(8)} {balance.coin}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-gray-600 dark:text-gray-400">Transfer Balance:</span>
          <span className="font-mono">{transferBalance.toFixed(8)} {balance.coin}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-gray-600 dark:text-gray-400">Locked Balance:</span>
          <span className="font-mono text-orange-600">{locked.toFixed(8)} {balance.coin}</span>
        </div>
        {bonus > 0 && (
          <div className="flex justify-between">
            <span className="text-gray-600 dark:text-gray-400">Bonus:</span>
            <span className="font-mono text-green-600">{bonus.toFixed(8)} {balance.coin}</span>
          </div>
        )}
        <div className="border-t pt-2 mt-2">
          <div className="flex justify-between font-medium">
            <span>Total Available:</span>
            <span className="font-mono">{(walletBalance + transferBalance + bonus).toFixed(8)} {balance.coin}</span>
          </div>
        </div>
      </div>
    </div>
  );
};

const CoinDetail: React.FC = () => {
  const { coinId } = useParams<{ coinId: string }>();
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [iconUrl, setIconUrl] = useState<string>("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [currentPrice, setCurrentPrice] = useState<string>("");
  const [priceLoading, setPriceLoading] = useState(false);
  const [lastUpdate, setLastUpdate] = useState<number>(0);
  const eventUnsubscribeRef = useRef<(() => void) | null>(null);

  useEffect(() => {
    if (!coinId) return;

    (async () => {
      try {
        // Fetch icon
        const icons = await GetCoinIconURLs([coinId.toUpperCase()]);
        
        if (icons && icons.length > 0) {
          const isDark = document.documentElement.classList.contains('dark');
          const icon = icons[0];
          const chosen = isDark
            ? (icon.darkUrl ?? icon.iconUrl)
            : (icon.lightUrl ?? icon.iconUrl);
          setIconUrl(chosen ?? '');
        }

        // Fetch current price
        setPriceLoading(true);
        try {
          const price = await GetCurrentPrice(coinId.toLowerCase());
          setCurrentPrice(price);
          
          // Start WebSocket stream for real-time updates
          await StartPriceStream(coinId.toLowerCase());
          
          // Listen for price updates
          const eventName = `price-update-${coinId.toLowerCase()}`;
          
          const unsubscribe = EventsOn(eventName, (priceData: any) => {
            if (priceData && typeof priceData === 'object' && priceData.price) {
              setCurrentPrice(priceData.price);
              setLastUpdate(Date.now());
            } else if (typeof priceData === 'string') {
              setCurrentPrice(priceData);
              setLastUpdate(Date.now());
            }
          });
          
          // Store unsubscribe function for cleanup
          eventUnsubscribeRef.current = unsubscribe;
          
        } catch (priceError) {
          console.error('Failed to fetch price:', priceError);
          setCurrentPrice('Unavailable');
        } finally {
          setPriceLoading(false);
        }

      } catch (e: any) {
        setError(String(e));
      } finally {
        setLoading(false);
      }
    })();

    // Cleanup function to stop price stream when component unmounts
    return () => {
      if (coinId) {
        StopPriceStream(coinId.toLowerCase()).catch(console.error);
        
        // Cleanup event listener
        if (eventUnsubscribeRef.current) {
          eventUnsubscribeRef.current();
        }
      }
    };
  }, [coinId]);

  if (loading) return <div>{t('Loading...')}</div>;
  if (error) return <div className="text-red-500">{error}</div>;

  return (
    <div className="p-4">
      <button
        onClick={() => navigate(-1)}
        className="mb-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition-colors"
      >
        {t('Back')}
      </button>
      
      <div className="flex items-center mb-4">
        {iconUrl && (
          <img
            src={iconUrl}
            alt={`${coinId} icon`}
            className="w-16 h-16 mr-4"
          />
        )}
        <h1 className="text-2xl font-bold text-gray-800 dark:text-white">
          {coinId?.toUpperCase()}
        </h1>
      </div>
      
      <div className="mt-4">
        {priceLoading ? (
          <div className="text-gray-500">{t('Loading price...')}</div>
        ) : (
          <div>
            <div className="text-2xl font-bold text-green-600">
              {currentPrice === 'Unavailable' ? (
                <span className="text-gray-500">Price unavailable</span>
              ) : (
                `$${currentPrice} USDT`
              )}
            </div>
            {lastUpdate > 0 && (
              <div className="flex items-center mt-2 text-sm text-gray-500">
                <div className="w-2 h-2 bg-green-500 rounded-full mr-2 animate-pulse"></div>
                Live
              </div>
            )}
          </div>
        )}
      </div>

      {/* Balance Section */}
      <div className="mt-6">
        <CoinBalanceSection coinSymbol={coinId || ''} />
      </div>
    </div>
  );
};

export default CoinDetail;
