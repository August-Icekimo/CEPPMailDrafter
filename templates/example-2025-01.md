---
from: cepp-notify@example.com
to:
  - primary@example.com
cc:
  - manager@example.com
subject: "{{YEAR}}年{{MONTH}}月 CEPP 例行通知"
attachments:
  - ./attachments/monthly-report.txt
---

<h2>{{YEAR}} 年 {{MONTH}} 月例行通知</h2>

<ul>
{{#ITEMS}}
  <li>{{ITEM}}</li>
{{/ITEMS}}
</ul>

{{#IF_ATTACHMENT}}
<p>本月附件請參閱隨附檔案。</p>
{{/IF_ATTACHMENT}}
