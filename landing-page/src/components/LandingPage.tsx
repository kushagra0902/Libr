import React from 'react';
import { Header, Hero, TechArch, WhatIsLIBR } from './LandingPageSections';
import { HowItWorks, Community, Footer } from './LandingPageExtended';
import { TechModules, HowToUse, InstallGuide } from './AdditionalSections';

interface LandingPageProps {
  isDarkMode: boolean;
  toggleTheme: () => void;
}

const LandingPage: React.FC<LandingPageProps> = ({ isDarkMode, toggleTheme }) => {
  return (
    <>
      <Header isDark={isDarkMode} toggleTheme={toggleTheme} />
      <Hero />  
      <WhatIsLIBR/>
      <HowToUse />
      <InstallGuide/>
      <TechArch />
      <HowItWorks />
      <TechModules />
      <Community />
      <Footer />
    </>
  );
};

export default LandingPage;
