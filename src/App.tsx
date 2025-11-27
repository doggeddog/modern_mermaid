import Layout from './components/Layout';
import { LanguageProvider } from './contexts/LanguageContext';
import { DarkModeProvider } from './contexts/DarkModeContext';
import GoogleAnalytics from './components/GoogleAnalytics';
import CookieConsent from './components/CookieConsent';

function App() {
  return (
    <DarkModeProvider>
      <LanguageProvider>
        <GoogleAnalytics />
        <CookieConsent />
        <Layout />
      </LanguageProvider>
    </DarkModeProvider>
  );
}

export default App;
