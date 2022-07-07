package pkg

var HTML = `<!DOCTYPE html>
<html lang="en" xmlns="http://www.w3.org/1999/html">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
    <script type="application/javascript">
        function search() {
            companyName = document.getElementById("companyName").value;
            alert(companyName)
        }
        function search() {
            var xhr = new XMLHttpRequest();
            page = document.getElementById("page").value;
            xhr.open('POST', '/search?page='+page);
            xhr.setRequestHeader("Content-type","application/json; charset=utf-8");
            var request = {CompanyName:document.getElementById("companyName").value};
            xhr.send(JSON.stringify(request));
            xhr.onreadystatechange = function () {
                if (xhr.readyState == 4) {
                    if (xhr.status = 200) {
                        var data = JSON.parse(xhr.responseText); //json解析方法JSON.parse 或者 eval('('+xhr.responseText+')')
                        count = xhr.getResponseHeader("count");
                        pagesize = parseInt(count/30)+1;
                        //alert(count,pagesize);
                        document.getElementById("count").value = count;
                        document.getElementById("pageSize").value = pagesize;
                        document.getElementById("tableBody").innerHTML = "";
                        for(var o in data){
                            //alert(data[o].CompanyName)
                            document.getElementById("tableBody").innerHTML +=
                                "<tr>" +
                                "<td>"+data[o].ID+"</td>" +
                                "<td>"+data[o].SearchName+"</td>" +
                                "<td>"+data[o].CompanyName+"</td>" +
                                "<td>"+data[o].LegalPerson+"</td>" +
                                "<td>"+data[o].Number+"</td>" +
                                "<td>"+data[o].Email+"</td>" +
                                "<td>"+data[o].Url+"</td>" +
                                "<td>"+data[o].Address+"</td>" +
                                "<td>"+data[o].CompanyAbbreviation+"</td>" +
                                //"<td>"+data[o].Introduction+"</td>" +
                                "<td>"+data[o].OperatingStatus+"</td>" +
                                "<td>"+data[o].DateOfEstablishment+"</td>" +
                                "<td>"+data[o].RegisteredCapital+"</td>" +
                                "<td>"+data[o].ApprovalDate+"</td>" +
                                "<td>"+data[o].Industry+"</td>" +
                                "<td>"+data[o].StaffSize+"</td>" +
                                "<td>"+data[o].NumberOfParticipants+"</td>" +
                                "<td>"+data[o].Area+"</td>" +
                                "<td>"+data[o].TotalMarketCapitalization+"</td>" +
                                "<td>"+data[o].MarketCapitalization+"</td>" +
                                "<td>"+data[o].OperatingIncome+"</td>" +
                                "<td>"+data[o].NetProfit+"</td>" +
                                "<td>"+data[o].BusinessArea+"</td>" +
                                "<td>"+data[o].NetInterestRate+"</td>" +
                                "<td>"+data[o].FinancingRounds+"</td>" +
                                //"<td>"+data[o].FinancingTime+"</td>" +
                                //"<td>"+data[o].FinancingRound+"</td>" +
                                //"<td>"+data[o].FinancingAmount+"</td>" +
                                //"<td>"+data[o].InvestmentAgency+"</td>" +
                                //"<td>"+data[o].Members+"</td>" +
                                //"<td>"+data[o].CompatProducts+"</td>" +
                                /*"<td>"+'<input type="button" value="查看详情" onclick=mainMembers("'+data[o].CompanyName+'") /></td>' +
                                "<td>"+'<input type="button" value="查看详情" onclick=competitors("'+data[o].CompanyName+'") /></td>' +*/
                                "</tr>"
                        }
                    }
                }
            }
        }
        function mainMembers(companyName) {
            var xhr = new XMLHttpRequest();
            xhr.open('POST', '/mainMembers');
            xhr.setRequestHeader("Content-type","application/json; charset=utf-8");
            var request = {CompanyName:companyName};
            xhr.send(JSON.stringify(request));
            xhr.onreadystatechange = function () {
                if (xhr.readyState == 4) {
                    if (xhr.status = 200) {
                        var data = JSON.parse(xhr.responseText); //json解析方法JSON.parse 或者 eval('('+xhr.responseText+')')
                        var value = ""
                        for(var o in data) {
                            value += data[o].Name +":"+data[o].Position+","
                        }
                        alert(value)
                    }
                }
            }
        }
        function competitors(companyName) {
            var xhr = new XMLHttpRequest();
            xhr.open('POST', '/competitors');
            xhr.setRequestHeader("Content-type","application/json; charset=utf-8");
            var request = {CompanyName:companyName};
            xhr.send(JSON.stringify(request));
            xhr.onreadystatechange = function () {
                if (xhr.readyState == 4) {
                    if (xhr.status = 200) {
                        var data = JSON.parse(xhr.responseText); //json解析方法JSON.parse 或者 eval('('+xhr.responseText+')')
                        var value = "";
                        for(var o in data) {
                            value += data[o].CompatName +","
                        }
                        alert(value)
                    }
                }
            }
        }
        function nextPage() {
            pageSize = parseInt(document.getElementById("pageSize").value);
            page = parseInt(document.getElementById("page").value);
            if(page < pageSize){
                page ++;
                document.getElementById("page").value = page;
                search()
            }
        }
        function prePage() {
            page = parseInt(document.getElementById("page").value);
            if(page > 1){
                page --;
                document.getElementById("page").value = page;
                search()
            }
        }
        function startSpider() {
            var xhr = new XMLHttpRequest();
            xhr.open('POST', '/startSpider');
            xhr.setRequestHeader("Content-type","application/json; charset=utf-8");
            var request = {CompanyName:companyName};
            xhr.send(JSON.stringify(request));
            xhr.onreadystatechange = function () {
                if (xhr.readyState == 4) {
                    if (xhr.status = 200) {
                        var data = JSON.parse(xhr.responseText); //json解析方法JSON.parse 或者 eval('('+xhr.responseText+')')
                        var value = "";
                        for(var o in data) {
                            value += data[o].CompatName +","
                        }
                        alert(value)
                    }
                }
            }
        }
        function exportAll() {
            var xhr = new XMLHttpRequest();
            xhr.open('POST', '/exportAll');
            xhr.setRequestHeader("Content-type","application/json; charset=utf-8");
            var request = {CompanyName:companyName};
            xhr.send(JSON.stringify(request));
            xhr.onreadystatechange = function () {
                if (xhr.readyState == 4) {
                    if (xhr.status = 200) {
                        alert("导出完成")
                    }
                }
            }
        }
    </script>
</head>
<body>
</br>
第<input style="width: 30px" disabled="disabled" value="1" id="page" />页&nbsp;&nbsp;共<input style="width: 30px" disabled="disabled" value="1" id="pageSize" />页&nbsp;&nbsp;共<input style="width: 30px" disabled="disabled"  value="1" id="count" />条
&nbsp;&nbsp;企业名称：<input type="text" name="companyName" id="companyName"/>&nbsp;&nbsp;<button onclick="search()">搜索</button>&nbsp;&nbsp;<button onclick="prePage()">上一页</button>&nbsp;&nbsp;<button onclick="nextPage()">下一页</button>
&nbsp;&nbsp;<button onclick="startSpider()">开始爬虫</button>
&nbsp;&nbsp;<button onclick="exportAll()">导出全部</button>
</br></br></br>
<table class="gridtable">
    <thead>
    <tr>
        <td>ID</td><td>搜索公司名称</td><td>公司名称</td><td>企业法人</td><td>联系电话</td><td>邮箱</td><td>官网地址</td><td>地址</td><td>品牌名称</td><td>经营状态</td><td>成立日期</td><td>注册资本</td><td>核准日期</td><td>所属行业</td><td>人员规模</td><td>参保人数</td><td>所属地区</td><td>总市值</td>
        <td>流通市值</td><td>营业收入</td><td>净利润</td><td>营业区间</td><td>净利率</td><td>融资伦次</td>
        <!--
        <td>融资时间</td><td>融资轮次</td><td>融资金额</td><td>投资机构</td>
        <td>主要人员</td><td>竞争企业</td>
        -->
    </tr>
    </thead>
    <tbody id="tableBody">
    </tbody>
</table>
</body>
<style>
    table.gridtable {
        font-family: verdana,arial,sans-serif;
        font-size:11px;
        color:#333333;
        border-width: 1px;
        border-color: #666666;
        border-collapse: collapse;
    }
    table.gridtable th {
        border-width: 1px;
        padding: 8px;
        border-style: solid;
        border-color: #666666;
        background-color: #dedede;
    }
    table.gridtable td {
        border-width: 1px;
        padding: 8px;
        border-style: solid;
        border-color: #666666;
        background-color: #ffffff;
        word-wrap:break-word; /*允许长单词换行到下一行*/
        word-break:break-all; /*这个参数根据需要来决定要不要*/
        overflow: hidden;/*这个参数根据需要来决定要不要*/
    }
</style>
</html>`
