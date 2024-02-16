import { Routes } from '@angular/router';
import { BtpComponent } from './btp/btp.component';
import { BotComponent } from './bot/bot.component';
import { HomepageComponent } from './homepage/homepage.component';

export const routes: Routes = [
    { path: 'btp', component: BtpComponent },
    { path: 'bot', component: BotComponent },
    { path: '', component: HomepageComponent },
    // { path: 'second-component', component: SecondComponent },
];
