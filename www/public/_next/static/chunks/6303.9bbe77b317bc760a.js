(self.webpackChunk_N_E=self.webpackChunk_N_E||[]).push([[6303],{2901:function(e,t,r){"use strict";r.r(t),r.d(t,{Artifacts:function(){return b},ArtifactsShareButton:function(){return y},HTMLPreview:function(){return C}});var n=r(57437),s=r(2265),i=r(59208),a=r(35499),o=r(84065),l=r(12637),c=r(84193),u=r(72495),d=r(59566),h=r(49111),f=r(98829),m=r(88587),x=r(65878),v=r(92944),g=r(75591),j=r(77120),p=r(99164),w=r.n(p);let C=(0,s.forwardRef)(function(e,t){let r=(0,s.useRef)(null),[i,a]=(0,s.useState)((0,o.x0)()),[l,c]=(0,s.useState)(600),[u,d]=(0,s.useState)("");(0,s.useEffect)(()=>{let e=e=>{let{id:t,height:r,title:n}=e.data;d(n),t==i&&c(r)};return window.addEventListener("message",e),()=>{window.removeEventListener("message",e)}},[i]),(0,s.useImperativeHandle)(t,()=>({reload:()=>{a((0,o.x0)())}}));let h=(0,s.useMemo)(()=>{if(!e.autoHeight)return e.height||600;if("string"==typeof e.height)return e.height;let t=e.height||600;return l+40>t?t:l+40},[e.autoHeight,e.height,l]),f=(0,s.useMemo)(()=>{let t='<script>window.addEventListener("DOMContentLoaded", () => new ResizeObserver((entries) => parent.postMessage({id: \''.concat(i,"', height: entries[0].target.clientHeight}, '*')).observe(document.body))</script>");return e.code.includes("<!DOCTYPE html>")&&e.code.replace("<!DOCTYPE html>","<!DOCTYPE html>"+t),t+e.code},[e.code,i]);return(0,n.jsx)("iframe",{className:w()["artifacts-iframe"],ref:r,sandbox:"allow-forms allow-modals allow-scripts",style:{height:h},srcDoc:f,onLoad:()=>{(null==e?void 0:e.onLoad)&&e.onLoad(u)}},i)});function y(e){let{getCode:t,id:r,style:i,fileName:o}=e,[d,f]=(0,s.useState)(!1),[j,p]=(0,s.useState)(r),[w,C]=(0,s.useState)(!1),y=(0,s.useMemo)(()=>[location.origin,"#",g.y$.Artifacts,"/",j].join(""),[j]),b=e=>r?Promise.resolve({id:r}):fetch(g.L.Artifacts,{method:"POST",body:e}).then(e=>e.json()).then(e=>{let{id:t}=e;if(t)return{id:t};throw Error()}).catch(e=>{(0,x.CF)(m.ZP.Export.Artifacts.Error)});return(0,n.jsxs)(n.Fragment,{children:[(0,n.jsx)("div",{className:"window-action-button",style:i,children:(0,n.jsx)(a.h,{icon:d?(0,n.jsx)(h.Z,{}):(0,n.jsx)(l.Z,{}),bordered:!0,title:m.ZP.Export.Artifacts.Title,onClick:()=>{d||(f(!0),b(t()).then(e=>{(null==e?void 0:e.id)&&(C(!0),p(null==e?void 0:e.id))}).finally(()=>f(!1)))}})}),w&&(0,n.jsx)("div",{className:"modal-mask",children:(0,n.jsx)(x.u_,{title:m.ZP.Export.Artifacts.Title,onClose:()=>C(!1),actions:[(0,n.jsx)(a.h,{icon:(0,n.jsx)(u.Z,{}),bordered:!0,text:m.ZP.Export.Download,onClick:()=>{(0,v.CP)(t(),"".concat(o||j,".html")).then(()=>C(!1))}},"download"),(0,n.jsx)(a.h,{icon:(0,n.jsx)(c.Z,{}),bordered:!0,text:m.ZP.Chat.Actions.Copy,onClick:()=>{(0,v.vQ)(y).then(()=>C(!1))}},"copy")],children:(0,n.jsx)("div",{children:(0,n.jsx)("a",{target:"_blank",href:y,children:y})})})})]})}function b(){let{id:e}=(0,i.UO)(),[t,r]=(0,s.useState)(""),[o,l]=(0,s.useState)(!0),[c,u]=(0,s.useState)(""),h=(0,s.useRef)(null);return(0,s.useEffect)(()=>{e&&fetch("".concat(g.L.Artifacts,"?id=").concat(e)).then(e=>{if(e.status>300)throw Error("can not get content");return e}).then(e=>e.text()).then(r).catch(e=>{(0,x.CF)(m.ZP.Export.Artifacts.Error)})},[e]),(0,n.jsxs)("div",{className:w().artifacts,children:[(0,n.jsxs)("div",{className:w()["artifacts-header"],children:[(0,n.jsx)("a",{href:g.Bv,target:"_blank",rel:"noopener noreferrer",children:(0,n.jsx)(a.h,{bordered:!0,icon:(0,n.jsx)(d.Z,{}),shadow:!0})}),(0,n.jsx)(a.h,{bordered:!0,style:{marginLeft:20},icon:(0,n.jsx)(f.Z,{}),shadow:!0,onClick:()=>{var e;return null===(e=h.current)||void 0===e?void 0:e.reload()}}),(0,n.jsx)("div",{className:w()["artifacts-title"],children:"NeatChat Artifacts"}),(0,n.jsx)(y,{id:e,getCode:()=>t,fileName:c})]}),(0,n.jsxs)("div",{className:w()["artifacts-content"],children:[o&&(0,n.jsx)(j.Loading,{}),t&&(0,n.jsx)(C,{code:t,ref:h,autoHeight:!1,height:"100%",onLoad:e=>{u(e),l(!1)}})]})]})}},42896:function(e,t,r){"use strict";r.r(t),r.d(t,{Markdown:function(){return N},MarkdownContent:function(){return P},Mermaid:function(){return E},PreCode:function(){return Z}});var n=r(57437),s=r(40145);r(68128);var i=r(28105),a=r(32794),o=r(50170),l=r(11303),c=r(74014),u=r(49264),d=r(2265),h=r(92944),f=r(51893),m=r(88587),x=r(24053),v=r(98829),g=r(38648),j=r(65878),p=r(2901),w=r(62425),C=r(35499),y=r(3289),b=r(75504);function k(e){return(0,n.jsx)("details",{children:e.children})}function S(e){return(0,n.jsx)("summary",{children:e.children})}function E(e){let t=(0,d.useRef)(null),[r,s]=(0,d.useState)(!1);return((0,d.useEffect)(()=>{e.code&&t.current&&f.L.run({nodes:[t.current],suppressErrors:!0}).catch(e=>{s(!0),console.error("[Mermaid] ",e.message)})},[e.code]),r)?null:(0,n.jsx)("div",{className:(0,b.Z)("no-dark","mermaid"),style:{cursor:"pointer",overflow:"auto"},ref:t,onClick:()=>(function(){var e;let r=null===(e=t.current)||void 0===e?void 0:e.querySelector("svg");if(!r)return;let n=new XMLSerializer().serializeToString(r),s=new Blob([n],{type:"image/svg+xml"});(0,j.vi)(URL.createObjectURL(s))})(),children:e.code})}function Z(e){var t;let r=(0,d.useRef)(null),s=(0,d.useRef)(null),[i,a]=(0,d.useState)(""),[o,l]=(0,d.useState)(""),{height:c}=(0,h.iP)(),u=(0,w.aK)().currentSession(),f=(0,g.y1)(()=>{var e;if(!r.current)return;let t=r.current.querySelector("code.language-mermaid");t&&a(t.innerText);let n=r.current.querySelector("code.language-html"),s=null===(e=r.current.querySelector("code"))||void 0===e?void 0:e.innerText;n?l(n.innerText):((null==s?void 0:s.startsWith("<!DOCTYPE"))||(null==s?void 0:s.startsWith("<svg"))||(null==s?void 0:s.startsWith("<?xml")))&&l(s)},600),m=(0,y.MG)(),x=(null===(t=u.mask)||void 0===t?void 0:t.enableArtifacts)!==!1&&m.enableArtifacts;return(0,d.useEffect)(()=>{if(r.current){let e=r.current.querySelectorAll("code"),t=["","md","markdown","text","txt","plaintext","tex","latex"];e.forEach(e=>{let r=e.className.match(/language-(\w+)/),n=r?r[1]:"";t.includes(n)&&(e.style.whiteSpace="pre-wrap")}),setTimeout(f,1)}},[]),(0,n.jsxs)(n.Fragment,{children:[(0,n.jsxs)("pre",{ref:r,children:[(0,n.jsx)("span",{className:"copy-code-button",onClick:()=>{if(r.current){var e,t;(0,h.vQ)(null!==(t=null===(e=r.current.querySelector("code"))||void 0===e?void 0:e.innerText)&&void 0!==t?t:"")}}}),e.children]}),i.length>0&&(0,n.jsx)(E,{code:i},i),o.length>0&&x&&(0,n.jsxs)(j.IT,{className:"no-dark html",right:70,children:[(0,n.jsx)(p.ArtifactsShareButton,{style:{position:"absolute",right:20,top:10},getCode:()=>o}),(0,n.jsx)(C.h,{style:{position:"absolute",right:120,top:10},bordered:!0,icon:(0,n.jsx)(v.Z,{}),shadow:!0,onClick:()=>{var e;return null===(e=s.current)||void 0===e?void 0:e.reload()}}),(0,n.jsx)(p.HTMLPreview,{ref:s,code:o,autoHeight:!document.fullscreenElement,height:document.fullscreenElement?c:600})]})]})}function _(e){var t;let r=(0,w.aK)().currentSession(),s=(0,y.MG)(),i=(null===(t=r.mask)||void 0===t?void 0:t.enableCodeFold)!==!1&&s.enableCodeFold,a=(0,d.useRef)(null),[o,l]=(0,d.useState)(!0),[c,u]=(0,d.useState)(!1);return(0,d.useEffect)(()=>{a.current&&(u(a.current.scrollHeight>400),a.current.scrollTop=a.current.scrollHeight)},[e.children]),(0,n.jsxs)(n.Fragment,{children:[(0,n.jsx)("code",{className:(0,b.Z)(null==e?void 0:e.className),ref:a,style:{maxHeight:i&&o?"400px":"none",overflowY:"hidden"},children:e.children}),c&&i&&o?(0,n.jsx)("div",{className:(0,b.Z)("show-hide-button",{collapsed:o,expanded:!o}),children:(0,n.jsx)("button",{onClick:()=>{l(e=>!e)},children:m.ZP.NewChat.More})}):null]})}let P=d.memo(function(e){let t=(0,d.useMemo)(()=>{var t;return(t=function(e){if(e.startsWith("<think>")&&!e.includes("</think>")){let t=e.slice(7);return"<details>\n<summary>".concat(m.ZP.NewChat.Thinking,"</summary>\n\n").concat(t,"\n\n</details>")}return e.replace(/^<think>([\s\S]*?)<\/think>/,(e,t)=>"<details>\n<summary>".concat(m.ZP.NewChat.Think,"</summary>\n\n").concat(t,"\n\n</details>"))}(e.content.replace(/(```[\s\S]*?```|`.*?`)|\\\[([\s\S]*?[^\\])\\\]|\\\((.*?)\\\)/g,(e,t,r,n)=>t||(r?"$$".concat(r,"$$"):n?"$".concat(n,"$"):e)))).includes("```")?t:t.replace(/([`]*?)(\w*?)([\n\r]*?)(<!DOCTYPE html>)/g,(e,t,r,n,s)=>t?e:"\n```html\n"+s).replace(/(<\/body>)([\r\n\s]*?)(<\/html>)([\n\r]*)([`]*)([\n\r]*?)/g,(e,t,r,n,s,i)=>i?e:t+r+n+"\n```\n")},[e.content]);return(0,n.jsx)(s.D,{remarkPlugins:[i.Z,l.Z,a.Z],rehypePlugins:[c.Z,o.Z,[u.Z,{detect:!1,ignoreMissing:!0}]],components:{pre:Z,code:_,p:e=>(0,n.jsx)("p",{...e,dir:"auto"}),a:e=>{var t;let r=e.href||"";if(/\.(aac|mp3|opus|wav)$/.test(r))return(0,n.jsx)("figure",{children:(0,n.jsx)("audio",{controls:!0,src:r})});if(/\.(3gp|3g2|webm|ogv|mpeg|mp4|avi)$/.test(r))return(0,n.jsx)("video",{controls:!0,width:"99.9%",children:(0,n.jsx)("source",{src:r})});let s=/^\/#/i.test(r)?"_self":null!==(t=e.target)&&void 0!==t?t:"_blank";return(0,n.jsx)("a",{...e,target:s})},details:k,summary:S},children:t})});function N(e){var t;let r=(0,d.useRef)(null);return(0,n.jsx)("div",{className:"markdown-body",style:{fontSize:"".concat(null!==(t=e.fontSize)&&void 0!==t?t:14,"px"),fontFamily:e.fontFamily||"inherit"},ref:r,onContextMenu:e.onContextMenu,onDoubleClickCapture:e.onDoubleClickCapture,dir:"auto",children:e.loading?(0,n.jsx)(x.Z,{}):(0,n.jsx)(P,{content:e.content})})}},99164:function(e){e.exports={artifacts:"artifacts_artifacts__J06vB","artifacts-header":"artifacts_artifacts-header__s7Cdi","artifacts-title":"artifacts_artifacts-title__UXZs9","artifacts-content":"artifacts_artifacts-content__3pFba","artifacts-iframe":"artifacts_artifacts-iframe__mjsdx"}}}]);