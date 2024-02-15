import { Routes } from '@angular/router';
import { BtpComponent } from './btp/btp.component';
import { BotComponent } from './bot/bot.component';

export const routes: Routes = [
    { path: 'btp', component: BtpComponent },
    { path: 'bot', component: BotComponent },
    // { path: 'second-component', component: SecondComponent },
];
