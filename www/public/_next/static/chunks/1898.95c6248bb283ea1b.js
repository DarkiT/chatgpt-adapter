"use strict";(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[1898],{41898:function(e,l,t){t.d(l,{diagram:function(){return g}});var a=t(63627),n=t(36783),i=t(77832),o=t(51893),s=t(15883);t(89539),t(99824),t(68985),t(91605),t(23780);let d=e=>o.e.sanitizeText(e,(0,o.c)()),r={dividerMargin:10,padding:5,textHeight:10,curve:void 0},c=function(e,l,t,a){let n=Object.keys(e);o.l.info("keys:",n),o.l.info(e),n.forEach(function(n){var i,s,r;let c=e[n],b={shape:"rect",id:c.id,domId:c.domId,labelText:d(c.id),labelStyle:"",style:"fill: none; stroke: black",padding:null!==(r=null==(i=(0,o.c)().flowchart)?void 0:i.padding)&&void 0!==r?r:null==(s=(0,o.c)().class)?void 0:s.padding};l.setNode(c.id,b),u(c.classes,l,t,a,c.id),o.l.info("setNode",b)})},u=function(e,l,t,a,n){let i=Object.keys(e);o.l.info("keys:",i),o.l.info(e),i.filter(l=>e[l].parent==n).forEach(function(t){var i,s,r,c;let u=e[t],b=u.cssClasses.join(" "),p={labelStyle:"",style:""},y=null!==(r=u.label)&&void 0!==r?r:u.id,f={labelStyle:p.labelStyle,shape:"class_box",labelText:d(y),classData:u,rx:0,ry:0,class:b,style:p.style,id:u.id,domId:u.domId,tooltip:a.db.getTooltip(u.id,n)||"",haveCallback:u.haveCallback,link:u.link,width:"group"===u.type?500:void 0,type:u.type,padding:null!==(c=null==(i=(0,o.c)().flowchart)?void 0:i.padding)&&void 0!==c?c:null==(s=(0,o.c)().class)?void 0:s.padding};l.setNode(u.id,f),n&&l.setParent(u.id,n),o.l.info("setNode",f)})},b=function(e,l,t,a){o.l.info(e),e.forEach(function(e,i){var s,c,u;let b={labelStyle:"",style:""},p=e.text,y={labelStyle:b.labelStyle,shape:"note",labelText:d(p),noteData:e,rx:0,ry:0,class:"",style:b.style,id:e.id,domId:e.id,tooltip:"",type:"note",padding:null!==(u=null==(s=(0,o.c)().flowchart)?void 0:s.padding)&&void 0!==u?u:null==(c=(0,o.c)().class)?void 0:c.padding};if(l.setNode(e.id,y),o.l.info("setNode",y),!e.class||!(e.class in a))return;let f=t+i,g={id:"edgeNote".concat(f),classes:"relation",pattern:"dotted",arrowhead:"none",startLabelRight:"",endLabelLeft:"",arrowTypeStart:"none",arrowTypeEnd:"none",style:"fill:none",labelStyle:"",curve:(0,o.n)(r.curve,n.c_6)};l.setEdge(e.id,e.class,g,f)})},p=function(e,l){let t=(0,o.c)().flowchart,a=0;e.forEach(function(e){var i,s;a++;let d={classes:"relation",pattern:1==e.relation.lineType?"dashed":"solid",id:"id"+a,arrowhead:"arrow_open"===e.type?"none":"normal",startLabelRight:"none"===e.relationTitle1?"":e.relationTitle1,endLabelLeft:"none"===e.relationTitle2?"":e.relationTitle2,arrowTypeStart:f(e.relation.type1),arrowTypeEnd:f(e.relation.type2),style:"fill:none",labelStyle:"",curve:(0,o.n)(null==t?void 0:t.curve,n.c_6)};if(o.l.info(d,e),void 0!==e.style){let l=(0,o.k)(e.style);d.style=l.style,d.labelStyle=l.labelStyle}e.text=e.title,void 0===e.text?void 0!==e.style&&(d.arrowheadStyle="fill: #333"):(d.arrowheadStyle="fill: #333",d.labelpos="c",(null!==(s=null==(i=(0,o.c)().flowchart)?void 0:i.htmlLabels)&&void 0!==s?s:(0,o.c)().htmlLabels)?(d.labelType="html",d.label='<span class="edgeLabel">'+e.text+"</span>"):(d.labelType="text",d.label=e.text.replace(o.e.lineBreakRegex,"\n"),void 0===e.style&&(d.style=d.style||"stroke: #333; stroke-width: 1.5px;fill:none"),d.labelStyle=d.labelStyle.replace("color:","fill:"))),l.setEdge(e.id1,e.id2,d,a)})},y=async function(e,l,t,a){var d,r,y,f;let g;o.l.info("Drawing class - ",l);let h=null!==(d=(0,o.c)().flowchart)&&void 0!==d?d:(0,o.c)().class,v=(0,o.c)().securityLevel;o.l.info("config:",h);let w=null!==(r=null==h?void 0:h.nodeSpacing)&&void 0!==r?r:50,k=null!==(y=null==h?void 0:h.rankSpacing)&&void 0!==y?y:50,x=new i.k({multigraph:!0,compound:!0}).setGraph({rankdir:a.db.getDirection(),nodesep:w,ranksep:k,marginx:8,marginy:8}).setDefaultEdgeLabel(function(){return{}}),m=a.db.getNamespaces(),S=a.db.getClasses(),T=a.db.getRelations(),L=a.db.getNotes();o.l.info(T),c(m,x,l,a),u(S,x,l,a),p(T,x),b(L,x,T.length+1,S),"sandbox"===v&&(g=(0,n.Ys)("#i"+l));let E="sandbox"===v?(0,n.Ys)(g.nodes()[0].contentDocument.body):(0,n.Ys)("body"),N=E.select('[id="'.concat(l,'"]')),D=E.select("#"+l+" g");if(await (0,s.r)(D,x,["aggregation","extension","composition","dependency","lollipop"],"classDiagram",l),o.u.insertTitle(N,"classTitleText",null!==(f=null==h?void 0:h.titleTopMargin)&&void 0!==f?f:5,a.db.getDiagramTitle()),(0,o.o)(x,N,null==h?void 0:h.diagramPadding,null==h?void 0:h.useMaxWidth),!(null==h?void 0:h.htmlLabels)){let e="sandbox"===v?g.nodes()[0].contentDocument:document;for(let t of e.querySelectorAll('[id="'+l+'"] .edgeLabel .label')){let l=t.getBBox(),a=e.createElementNS("http://www.w3.org/2000/svg","rect");a.setAttribute("rx",0),a.setAttribute("ry",0),a.setAttribute("width",l.width),a.setAttribute("height",l.height),t.insertBefore(a,t.firstChild)}}};function f(e){let l;switch(e){case 0:l="aggregation";break;case 1:l="extension";break;case 2:l="composition";break;case 3:l="dependency";break;case 4:l="lollipop";break;default:l="none"}return l}let g={parser:a.p,db:a.d,renderer:{setConf:function(e){r={...r,...e}},draw:y},styles:a.s,init:e=>{e.class||(e.class={}),e.class.arrowMarkerAbsolute=e.arrowMarkerAbsolute,a.d.clear()}}}}]);