import { Component } from '@angular/core';
import { Router } from '@angular/router';
import {MatIconModule} from '@angular/material/icon';

@Component({
  selector: 'app-homepage',
  standalone: true,
  imports: [MatIconModule],
  template: `
  <div class="w-full h-full bg-[#1f2836]">
    <div 
      
      class="absolute z-50 w-full h-[90%] flex items-center justify-center">
      <div>
        <div class=" pl-14 text-white text-[100px] " style="text-shadow: 2px 0 #000, -2px 0 #000, 0 2px #000, 0 -2px #000, 1px 1px #000, -1px -1px #000, 1px -1px #000, -1px 1px #000;">Welcome in BTP Tracker</div>
        <div class="w-full flex justify-center items-center h-[300px] mt-28 text-white">
          <div class="p-3 h-full w-full flex justify-center border-r-[1px] border-white" style="text-shadow: 1px 0 #000, -1px 0 #000, 0 1px #000, 0 -1px #000, 0.1px 0.1px #000, -0.1px -0.1px #000, 0.1px -0.1px #000, -0.1px 0.1px #000;">
            <div class="w-full h-full text-center">
              <div class="font-bold text-2xl mb-3">BTP</div>
              <div class="max-w-[600px] break-words text-center text-[#ffffff]">
              I Btp, Buoni del Tesoro Poliennali, sono titoli a medio-lungo termine,
              con una cedola fissa pagata semestralmente.
L’asta è riservata agli intermediari istituzionali 
autorizzati ai sensi del decreto legislativo 24 febbraio 1998, n. 58
              </div>
              <div (click)="changeRoute('btp')" class="font-bold text-left w-[70%] m-auto h-[50px] flex cursor-pointer px-7 items-center border mt-52 rounded-md border-[#784fbe] bg-[#1f2836] justify-between"><div>Go to btp history</div> <mat-icon aria-hidden="false" fontIcon="keyboard_arrow_right"></mat-icon> </div>
            </div>
          </div>
          <div class="h-full w-full flex justify-center">
          <div class="p-3 h-full w-full flex justify-center">
          <div class="w-full h-full text-center">
              <div class="font-bold text-2xl mb-3">BOT</div>
              <div class="max-w-[500px] break-words text-center text-[#ffffff]">
              I Bot sono titoli a breve termine con scadenza non superiore ad un anno. La remunerazione, interamente determinata dallo scarto di emissione (dato dalla differenza tra il valore nominale ed il prezzo pagato), è considerata ai fini fiscali anticipata, in quanto la ritenuta per gli investitori individuali si applica al momento della sottoscrizione.
              </div>
              
              <div (click)="changeRoute('bot')" class="font-bold text-left w-[70%] m-auto h-[50px] flex cursor-pointer px-7 items-center border mt-40 rounded-md border-[#784fbe] bg-[#1f2836] justify-between"><div>Go to bot history</div> <mat-icon aria-hidden="false" fontIcon="keyboard_arrow_right"></mat-icon> </div>
            </div>
          </div>
          </div>
        </div>
      </div>
    </div>




<div class="z-10 background">
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
   <span></span>
</div>
  </div>
  `,
  styleUrl: './homepage.component.css'
})
export class HomepageComponent {
  constructor(private router: Router) { }
  changeRoute (route:string) {
    this.router.navigate([`/${route}`]);
  }
}
