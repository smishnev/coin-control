import React, { useEffect, useState, useRef } from "react";
import { useTranslation } from "react-i18next";
import { useParams, useNavigate } from "react-router-dom";
import { GetCoinIconURLs, GetCurrentPrice, StartPriceStream, StopPriceStream } from "../../wailsjs/go/main/App";
import { EventsOn } from "../../wailsjs/runtime/runtime";

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
    </div>
  );
};

export default CoinDetail;
