import { HttpClient, HttpClientModule } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
/* Core Grid CSS */
import 'ag-grid-community/styles/ag-grid.css';
/* Quartz Theme Specific CSS */
import 'ag-grid-community/styles/ag-theme-quartz.css';
import { Router, RouterOutlet } from '@angular/router';
import { AgGridAngular } from 'ag-grid-angular';
import { ColDef } from 'ag-grid-community';
import { NgxChartsModule } from '@swimlane/ngx-charts';
import { ScaleType } from '@swimlane/ngx-charts';
import { BtpComponent } from './btp/btp.component';
import { BrowserModule } from '@angular/platform-browser';
// @ts-nocheck
@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, BtpComponent],
  template: `
  <div>
    <div class="h-[50px] flex justify-between items-center px-5 bg-[#11161e]">
    <div class="text-white font-black text-lg" (click)="changeRoute('/')">BTP Tracker</div>  
    <div class="flex justify-center items-center gap-4">
        
  <div class="border rounded-md py-1 px-4 text-white border-[#784fbe] bg-[#1f2836]" routerLink="/btp" (click)="changeRoute('btp')" routerLinkActive="active">BTP</div>
  <div class="border rounded-md py-1 px-4 text-white border-[#784fbe] bg-[#1f2836]" routerLink="/bot" (click)="changeRoute('bot')" routerLinkActive="active">BOT</div>
        
      </div>
    </div>
    <div class="absolute w-full" style="height: calc(100% - 50px);"><router-outlet/></div>
  </div>
  `,
  styleUrl: './app.component.css'
})
export class AppComponent {
  constructor(private router: Router) { }
  changeRoute (route:string) {
    this.router.navigate([`/${route}`]);
  }
 }