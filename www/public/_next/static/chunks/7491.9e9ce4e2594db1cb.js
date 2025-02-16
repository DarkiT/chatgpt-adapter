"use strict";(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[7491],{17491:function(e,t,o){o.d(t,{diagram:function(){return _}});var a=o(51421),s=o(77832),n=o(36783),i=o(51893),c=o(15883);o(89539),o(99824),o(68985),o(91605),o(23780);let r="rect",l="rectWithTitle",d="statediagram",p="".concat(d,"-").concat("state"),g="transition",b="".concat(g," ").concat("note-edge"),h="".concat(d,"-").concat("note"),u="".concat(d,"-").concat("cluster"),y="".concat(d,"-").concat("cluster-alt"),f="parent",w="note",x="----",m="".concat(x).concat(w),T="".concat(x).concat(f),S="fill:none",k="fill: #333",v="text",D="normal",A={},B=0;function E(){let e=arguments.length>0&&void 0!==arguments[0]?arguments[0]:"",t=arguments.length>1&&void 0!==arguments[1]?arguments[1]:0,o=arguments.length>2&&void 0!==arguments[2]?arguments[2]:"",a=arguments.length>3&&void 0!==arguments[3]?arguments[3]:x,s=null!==o&&o.length>0?"".concat(a).concat(o):"";return"".concat("state","-").concat(e).concat(s,"-").concat(t)}let N=(e,t,o,s,n,c)=>{var d;let g=o.id,x=null==(d=s[g])?"":d.classes?d.classes.join(" "):"";if("root"!==g){let t=r;!0===o.start&&(t="start"),!1===o.start&&(t="end"),o.type!==a.D&&(t=o.type),A[g]||(A[g]={id:g,shape:t,description:i.e.sanitizeText(g,(0,i.c)()),classes:"".concat(x," ").concat(p)});let s=A[g];o.description&&(Array.isArray(s.description)?(s.shape=l,s.description.push(o.description)):s.description.length>0?(s.shape=l,s.description===g?s.description=[o.description]:s.description=[s.description,o.description]):(s.shape=r,s.description=o.description),s.description=i.e.sanitizeTextOrArray(s.description,(0,i.c)())),1===s.description.length&&s.shape===l&&(s.shape=r),!s.type&&o.doc&&(i.l.info("Setting cluster for ",g,R(o)),s.type="group",s.dir=R(o),s.shape=o.type===a.a?"divider":"roundedWithTitle",s.classes=s.classes+" "+u+" "+(c?y:""));let n={labelStyle:"",shape:s.shape,labelText:s.description,classes:s.classes,style:"",id:g,dir:s.dir,domId:E(g,B),type:s.type,padding:15};if(n.centerLabel=!0,o.note){let t={labelStyle:"",shape:"note",labelText:o.note.text,classes:h,style:"",id:g+m+"-"+B,domId:E(g,B,w),type:s.type,padding:15},a={labelStyle:"",shape:"noteGroup",labelText:o.note.text,classes:s.classes,style:"",id:g+T,domId:E(g,B,f),type:"group",padding:0};B++;let i=g+T;e.setNode(i,a),e.setNode(t.id,t),e.setNode(g,n),e.setParent(g,i),e.setParent(t.id,i);let c=g,r=t.id;"left of"===o.note.position&&(c=t.id,r=g),e.setEdge(c,r,{arrowhead:"none",arrowType:"",style:S,labelStyle:"",classes:b,arrowheadStyle:k,labelpos:"c",labelType:v,thickness:D})}else e.setNode(g,n)}t&&"root"!==t.id&&(i.l.trace("Setting node ",g," to be child of its parent ",t.id),e.setParent(g,t.id)),o.doc&&(i.l.trace("Adding nodes children "),C(e,o,o.doc,s,n,!c))},C=(e,t,o,s,n,c)=>{i.l.trace("items",o),o.forEach(o=>{switch(o.stmt){case a.b:case a.D:N(e,t,o,s,n,c);break;case a.S:{N(e,t,o.state1,s,n,c),N(e,t,o.state2,s,n,c);let a={id:"edge"+B,arrowhead:"normal",arrowTypeEnd:"arrow_barb",style:S,labelStyle:"",label:i.e.sanitizeText(o.description,(0,i.c)()),arrowheadStyle:k,labelpos:"c",labelType:v,thickness:D,classes:g};e.setEdge(o.state1.id,o.state2.id,a,B),B++}}})},R=function(e){let t=arguments.length>1&&void 0!==arguments[1]?arguments[1]:a.c,o=t;if(e.doc)for(let t=0;t<e.doc.length;t++){let a=e.doc[t];"dir"===a.stmt&&(o=a.value)}return o},V=async function(e,t,o,a){let l;i.l.info("Drawing state diagram (v2)",t),A={},a.db.getDirection();let{securityLevel:p,state:g}=(0,i.c)(),b=g.nodeSpacing||50,h=g.rankSpacing||50;i.l.info(a.db.getRootDocV2()),a.db.extract(a.db.getRootDocV2()),i.l.info(a.db.getRootDocV2());let u=a.db.getStates(),y=new s.k({multigraph:!0,compound:!0}).setGraph({rankdir:R(a.db.getRootDocV2()),nodesep:b,ranksep:h,marginx:8,marginy:8}).setDefaultEdgeLabel(function(){return{}});N(y,void 0,a.db.getRootDocV2(),u,a.db,!0),"sandbox"===p&&(l=(0,n.Ys)("#i"+t));let f="sandbox"===p?(0,n.Ys)(l.nodes()[0].contentDocument.body):(0,n.Ys)("body"),w=f.select('[id="'.concat(t,'"]')),x=f.select("#"+t+" g");await (0,c.r)(x,y,["barb"],d,t),i.u.insertTitle(w,"statediagramTitleText",g.titleTopMargin,a.db.getDiagramTitle());let m=w.node().getBBox(),T=m.width+16,S=m.height+16;w.attr("class",d);let k=w.node().getBBox();(0,i.i)(w,S,T,g.useMaxWidth);let v="".concat(k.x-8," ").concat(k.y-8," ").concat(T," ").concat(S);for(let e of(i.l.debug("viewBox ".concat(v)),w.attr("viewBox",v),document.querySelectorAll('[id="'+t+'"] .edgeLabel .label'))){let t=e.getBBox(),o=document.createElementNS("http://www.w3.org/2000/svg",r);o.setAttribute("rx",0),o.setAttribute("ry",0),o.setAttribute("width",t.width),o.setAttribute("height",t.height),e.insertBefore(o,e.firstChild)}},_={parser:a.p,db:a.d,renderer:{setConf:function(e){for(let t of Object.keys(e))e[t]},getClasses:function(e,t){return t.db.extract(t.db.getRootDocV2()),t.db.getClasses()},draw:V},styles:a.s,init:e=>{e.state||(e.state={}),e.state.arrowMarkerAbsolute=e.arrowMarkerAbsolute,a.d.clear()}}}}]);