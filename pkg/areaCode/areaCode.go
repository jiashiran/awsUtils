package areaCode

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"spider/pkg"
	"spider/pkg/orm"
	"strings"
	"time"
)

func Get() {
	html := `<tr>
                      <th height="32" width="27%"><font color="#000000" size="2">Countries and Regions</font></th>
                      <th align="middle" height="32" width="24%"><font color="#000000" size="2">国家或地区</font></th>
                      <th align="middle" height="32" width="20%"><b><font face="宋体" size="2">国际域名缩写</font></b></th>
                      <th align="middle" height="32" width="15%"><font color="#000000" size="2">电话代码 </font></th>
                      <th align="middle" height="32" width="14%"><font size="2">时差</font></th></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Angola</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">安哥拉</font></td>
                      <td align="middle" width="20%"><font size="2">AO</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">244</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%">Afghanistan</td>
                      <td align="middle" width="24%"><font size="2">阿富汗</font></td>
                      <td align="middle" width="20%"><font size="2">AF</font></td>
                      <td align="middle" width="15%"><font size="2">93</font></td>
                      <td align="middle" width="14%"><font size="2">0</font></td></tr>
                    <tr>
                      <td width="27%">Albania</td>
                      <td align="middle" width="24%"><font size="2">阿尔巴尼亚</font></td>
                      <td align="middle" width="20%"><font size="2">AL</font></td>
                      <td align="middle" width="15%"><font size="2">355</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Algeria</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">阿尔及利亚</font></td>
                      <td align="middle" width="20%"><font size="2">DZ</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">213</font></td>
                      <td align="middle" width="14%"><font size="2">-8</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Andorra</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">安道尔共和国</font></td>
                      <td align="middle" width="20%"><font size="2">AD</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">376</font></td>
                      <td align="middle" width="14%"><font size="2">-8</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Anguilla</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">安圭拉岛</font></td>
                      <td align="middle" width="20%"><font size="2">AI</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">1264</font></td>
                      <td align="middle" width="14%"><font size="2">-12</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Antigua and Barbuda</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">安提瓜和巴布达</font></td>
                      <td align="middle" width="20%"><font size="2">AG</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">1268</font></td>
                      <td align="middle" width="14%"><font size="2">-12</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Argentina</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">阿根廷</font></td>
                      <td align="middle" width="20%"><font size="2">AR</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">54</font></td>
                      <td align="middle" width="14%"><font size="2">-11</font></td></tr>
                    <tr>
                      <td width="27%">Armenia</td>
                      <td align="middle" width="24%"><font size="2">亚美尼亚</font></td>
                      <td align="middle" width="20%"><font size="2">AM</font></td>
                      <td align="middle" width="15%"><font size="2">374</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Ascension</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">阿森松</font></td>
                      <td align="middle" width="20%"><font size="2">&nbsp;</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">247</font></td>
                      <td align="middle" width="14%"><font size="2">-8</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Australia</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">澳大利亚</font></td>
                      <td align="middle" width="20%"><font size="2">AU</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">61</font></td>
                      <td align="middle" width="14%"><font size="2">+2</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Austria</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">奥地利</font></td>
                      <td align="middle" width="20%"><font size="2">AT</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">43</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%">Azerbaijan</td>
                      <td align="middle" width="24%"><font size="2">阿塞拜疆</font></td>
                      <td align="middle" width="20%"><font size="2">AZ</font></td>
                      <td align="middle" width="15%"><font size="2">994</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Bahamas</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">巴哈马</font></td>
                      <td align="middle" width="20%"><font size="2">BS</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">1242</font></td>
                      <td align="middle" width="14%"><font size="2">-13</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Bahrain</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">巴林</font></td>
                      <td align="middle" width="20%"><font size="2">BH</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">973</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Bangladesh</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">孟加拉国</font></td>
                      <td align="middle" width="20%"><font size="2">BD</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">880</font></td>
                      <td align="middle" width="14%"><font size="2">-2</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Barbados</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">巴巴多斯</font></td>
                      <td align="middle" width="20%"><font size="2">BB</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">1246</font></td>
                      <td align="middle" width="14%"><font size="2">-12</font></td></tr>
                    <tr>
                      <td width="27%">Belarus</td>
                      <td align="middle" width="24%"><font size="2">白俄罗斯</font></td>
                      <td align="middle" width="20%"><font size="2">BY</font></td>
                      <td align="middle" width="15%"><font size="2">375</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Belgium</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">比利时</font></td>
                      <td align="middle" width="20%"><font size="2">BE</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">32</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Belize</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">伯利兹</font></td>
                      <td align="middle" width="20%"><font size="2">BZ</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">501</font></td>
                      <td align="middle" width="14%"><font size="2">-14</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Benin</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">贝宁</font></td>
                      <td align="middle" width="20%"><font size="2">BJ</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">229</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Bermuda</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">百慕大群岛</font></td>
                      <td align="middle" width="20%"><font size="2">BM</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">1441</font></td>
                      <td align="middle" width="14%"><font size="2">-12</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Bolivia</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">玻利维亚</font></td>
                      <td align="middle" width="20%"><font size="2">BO</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">591</font></td>
                      <td align="middle" width="14%"><font size="2">-12</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Botswana</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">博茨瓦纳</font></td>
                      <td align="middle" width="20%"><font size="2">BW</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">267</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Brazil</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">巴西</font></td>
                      <td align="middle" width="20%"><font size="2">BR</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">55</font></td>
                      <td align="middle" width="14%"><font size="2">-11</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Brunei</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">文莱</font></td>
                      <td align="middle" width="20%"><font size="2">BN</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">673</font></td>
                      <td align="middle" width="14%"><font size="2">0</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Bulgaria</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">保加利亚</font></td>
                      <td align="middle" width="20%"><font size="2">BG</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">359</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Burkina-faso</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">布基纳法索</font></td>
                      <td align="middle" width="20%"><font size="2">BF</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">226</font></td>
                      <td align="middle" width="14%"><font size="2">-8</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Burma</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">缅甸</font></td>
                      <td align="middle" width="20%"><font size="2">MM</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">95</font></td>
                      <td align="middle" width="14%"><font size="2">-1.3</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Burundi</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">布隆迪</font></td>
                      <td align="middle" width="20%"><font size="2">BI</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">257</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Cameroon</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">喀麦隆</font></td>
                      <td align="middle" width="20%"><font size="2">CM</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">237</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Canada</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">加拿大</font></td>
                      <td align="middle" width="20%"><font size="2">CA</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">1</font></td>
                      <td align="middle" width="14%"><font size="2">-13</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Cayman</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">开曼群岛</font></td>
                      <td align="middle" width="20%"><font size="2">&nbsp;</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">1345</font></td>
                      <td align="middle" width="14%"><font size="2">-13</font></td></tr>
                    <tr>
                      <td width="27%">Central African Republic</td>
                      <td align="middle" width="24%"><font size="2">中非共和国</font></td>
                      <td align="middle" width="20%"><font size="2">CF</font></td>
                      <td align="middle" width="15%"><font size="2">236</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%">Chad</td>
                      <td align="middle" width="24%"><font size="2">乍得</font></td>
                      <td align="middle" width="20%"><font size="2">TD</font></td>
                      <td align="middle" width="15%"><font size="2">235</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Chile</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">智利</font></td>
                      <td align="middle" width="20%"><font size="2">CL</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">56</font></td>
                      <td align="middle" width="14%"><font size="2">-13</font></td></tr>
                    <tr>
                      <td width="27%">China</td>
                      <td align="middle" width="24%"><font size="2">中国</font></td>
                      <td align="middle" width="20%"><font size="2">CN</font></td>
                      <td align="middle" width="15%"><font size="2">86</font></td>
                      <td align="middle" width="14%"><font size="2">0</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Colombia</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">哥伦比亚</font></td>
                      <td align="middle" width="20%"><font size="2">CO</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">57</font></td>
                      <td align="middle" width="14%"><font size="2">0</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Congo</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">刚果</font></td>
                      <td align="middle" width="20%"><font size="2">CG</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">242</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Cook</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">库克群岛</font></td>
                      <td align="middle" width="20%"><font size="2">CK</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">682</font></td>
                      <td align="middle" width="14%"><font size="2">-18.3</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Costa</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">哥斯达黎加</font></td>
                      <td align="middle" width="20%"><font size="2">CR</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">506</font></td>
                      <td align="middle" width="14%"><font size="2">-14</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Cuba</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">古巴</font></td>
                      <td align="middle" width="20%"><font size="2">CU</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">53</font></td>
                      <td align="middle" width="14%"><font size="2">-13</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Cyprus</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">塞浦路斯</font></td>
                      <td align="middle" width="20%"><font size="2">CY</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">357</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%">Czech Republic</td>
                      <td align="middle" width="24%"><font color="#000000" size="2">捷克</font></td>
                      <td align="middle" width="20%"><font size="2">CZ</font></td>
                      <td align="middle" width="15%"><font size="2">420</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Denmark</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">丹麦</font></td>
                      <td align="middle" width="20%"><font size="2">DK</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">45</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Djibouti</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">吉布提</font></td>
                      <td align="middle" width="20%"><font size="2">DJ</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">253</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Dominica</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">多米尼加共和国</font></td>
                      <td align="middle" width="20%"><font size="2">DO</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">1890</font></td>
                      <td align="middle" width="14%"><font size="2">-13</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Ecuador</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">厄瓜多尔</font></td>
                      <td align="middle" width="20%"><font size="2">EC</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">593</font></td>
                      <td align="middle" width="14%"><font size="2">-13</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Egypt</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">埃及</font></td>
                      <td align="middle" width="20%"><font size="2">EG</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">20</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">EL SALVADOR</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">萨尔瓦多</font></td>
                      <td align="middle" width="20%"><font size="2">SV</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">503</font></td>
                      <td align="middle" width="14%"><font size="2">-14</font></td></tr>
                    <tr>
                      <td height="17" width="27%">Estonia</td>
                      <td align="middle" height="17" width="24%"><font size="2">爱沙尼亚</font></td>
                      <td align="middle" height="17" width="20%"><font size="2">EE</font></td>
                      <td align="middle" height="17" width="15%"><font size="2">372</font></td>
                      <td align="middle" height="17" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Ethiopia</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">埃塞俄比亚</font></td>
                      <td align="middle" width="20%"><font size="2">ET</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">251</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Fiji</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">斐济</font></td>
                      <td align="middle" width="20%"><font size="2">FJ</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">679</font></td>
                      <td align="middle" width="14%"><font size="2">+4</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Finland</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">芬兰</font></td>
                      <td align="middle" width="20%"><font size="2">FI</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">358</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">France</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">法国</font></td>
                      <td align="middle" width="20%"><font size="2">FR</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">33</font></td>
                      <td align="middle" width="14%"><font size="2">-8</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">French 
                      Guiana</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">法属圭亚那</font></td>
                      <td align="middle" width="20%"><font size="2">GF</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">594</font></td>
                      <td align="middle" width="14%"><font size="2">-12</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Gabon</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">加蓬</font></td>
                      <td align="middle" width="20%"><font size="2">GA</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">241</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Gambia</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">冈比亚</font></td>
                      <td align="middle" width="20%"><font size="2">GM</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">220</font></td>
                      <td align="middle" width="14%"><font size="2">-8</font></td></tr>
                    <tr>
                      <td width="27%">Georgia </td>
                      <td align="middle" width="24%"><font size="2">格鲁吉亚</font></td>
                      <td align="middle" width="20%"><font size="2">GE</font></td>
                      <td align="middle" width="15%"><font size="2">995</font></td>
                      <td align="middle" width="14%"><font size="2">0</font></td></tr>
                    <tr>
                      <td width="27%">Germany </td>
                      <td align="middle" width="24%"><font size="2">德国</font></td>
                      <td align="middle" width="20%"><font size="2">DE</font></td>
                      <td align="middle" width="15%"><font size="2">49</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%">Ghana</td>
                      <td align="middle" width="24%"><font size="2">加纳</font></td>
                      <td align="middle" width="20%"><font size="2">GH</font></td>
                      <td align="middle" width="15%"><font size="2">233</font></td>
                      <td align="middle" width="14%"><font size="2">-8</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Gibraltar</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">直布罗陀</font></td>
                      <td align="middle" width="20%"><font size="2">GI</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">350</font></td>
                      <td align="middle" width="14%"><font size="2">-8</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Greece</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">希腊</font></td>
                      <td align="middle" width="20%"><font size="2">GR</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">30</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Grenada</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">格林纳达</font></td>
                      <td align="middle" width="20%"><font size="2">GD</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">1809</font></td>
                      <td align="middle" width="14%"><font size="2">-14</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Guam</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">关岛</font></td>
                      <td align="middle" width="20%"><font size="2">GU</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">1671</font></td>
                      <td align="middle" width="14%"><font size="2">+2</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Guatemala</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">危地马拉</font></td>
                      <td align="middle" width="20%"><font size="2">GT</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">502</font></td>
                      <td align="middle" width="14%"><font size="2">-14</font></td></tr>
                    <tr>
                      <td width="27%">Guinea</td>
                      <td align="middle" width="24%"><font size="2">几内亚</font></td>
                      <td align="middle" width="20%"><font size="2">GN</font></td>
                      <td align="middle" width="15%"><font size="2">224</font></td>
                      <td align="middle" width="14%"><font size="2">-8</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Guyana</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">圭亚那</font></td>
                      <td align="middle" width="20%"><font size="2">GY</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">592</font></td>
                      <td align="middle" width="14%"><font size="2">-11</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Haiti</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">海地</font></td>
                      <td align="middle" width="20%"><font size="2">HT</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">509</font></td>
                      <td align="middle" width="14%"><font size="2">-13</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Honduras</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">洪都拉斯</font></td>
                      <td align="middle" width="20%"><font size="2">HN</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">504</font></td>
                      <td align="middle" width="14%"><font size="2">-14</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Hongkong</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">香港</font></td>
                      <td align="middle" width="20%"><font size="2">HK</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">852</font></td>
                      <td align="middle" width="14%"><font size="2">0</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Hungary</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">匈牙利</font></td>
                      <td align="middle" width="20%"><font size="2">HU</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">36</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Iceland</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">冰岛</font></td>
                      <td align="middle" width="20%"><font size="2">IS</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">354</font></td>
                      <td align="middle" width="14%"><font size="2">-9</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">India</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">印度</font></td>
                      <td align="middle" width="20%"><font size="2">IN</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">91</font></td>
                      <td align="middle" width="14%"><font size="2">-2.3</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Indonesia</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">印度尼西亚</font></td>
                      <td align="middle" width="20%"><font size="2">ID</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">62</font></td>
                      <td align="middle" width="14%"><font size="2">-0.3</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Iran</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">伊朗</font></td>
                      <td align="middle" width="20%"><font size="2">IR</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">98</font></td>
                      <td align="middle" width="14%"><font size="2">-4.3</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Iraq</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">伊拉克</font></td>
                      <td align="middle" width="20%"><font size="2">IQ</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">964</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Ireland</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">爱尔兰</font></td>
                      <td align="middle" width="20%"><font size="2">IE</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">353</font></td>
                      <td align="middle" width="14%"><font size="2">-4.3</font></td></tr>
                    <tr>
                      <td width="27%">Israel</td>
                      <td align="middle" width="24%"><font size="2">以色列</font></td>
                      <td align="middle" width="20%"><font size="2">IL</font></td>
                      <td align="middle" width="15%"><font size="2">972</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Italy</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">意大利</font></td>
                      <td align="middle" width="20%"><font size="2">IT</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">39</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Ivory 
Coast</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">科特迪瓦</font></td>
                      <td align="middle" width="20%"><font size="2">&nbsp;</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">225</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Jamaica</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">牙买加</font></td>
                      <td align="middle" width="20%"><font size="2">JM</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">1876</font></td>
                      <td align="middle" width="14%"><font size="2">-12</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Japan</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">日本</font></td>
                      <td align="middle" width="20%"><font size="2">JP</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">81</font></td>
                      <td align="middle" width="14%"><font size="2">+1</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Jordan</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">约旦</font></td>
                      <td align="middle" width="20%"><font size="2">JO</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">962</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%">Kampuchea (Cambodia )</td>
                      <td align="middle" width="24%"><font size="2">柬埔寨</font></td>
                      <td align="middle" width="20%"><font size="2">KH</font></td>
                      <td align="middle" width="15%"><font size="2">855</font></td>
                      <td align="middle" width="14%"><font size="2">-1</font></td></tr>
                    <tr>
                      <td width="27%">Kazakstan</td>
                      <td align="middle" width="24%"><font size="2">哈萨克斯坦</font></td>
                      <td align="middle" width="20%"><font size="2">KZ</font></td>
                      <td align="middle" width="15%"><font size="2">327</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Kenya</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">肯尼亚</font></td>
                      <td align="middle" width="20%"><font size="2">KE</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">254</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%">Korea</td>
                      <td align="middle" width="24%"><font size="2">韩国</font></td>
                      <td align="middle" width="20%"><font size="2">KR</font></td>
                      <td align="middle" width="15%"><font size="2">82</font></td>
                      <td align="middle" width="14%"><font size="2">+1</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Kuwait</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">科威特</font></td>
                      <td align="middle" width="20%"><font size="2">KW</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">965</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%">Kyrgyzstan </td>
                      <td align="middle" width="24%"><font size="2">吉尔吉斯坦</font></td>
                      <td align="middle" width="20%"><font size="2">KG</font></td>
                      <td align="middle" width="15%"><font size="2">331</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%">Laos</td>
                      <td align="middle" width="24%"><font size="2">老挝</font></td>
                      <td align="middle" width="20%"><font size="2">LA</font></td>
                      <td align="middle" width="15%"><font size="2">856</font></td>
                      <td align="middle" width="14%"><font size="2">-1</font></td></tr>
                    <tr>
                      <td width="27%">Latvia </td>
                      <td align="middle" width="24%"><font size="2">拉脱维亚</font></td>
                      <td align="middle" width="20%"><font size="2">LV</font></td>
                      <td align="middle" width="15%"><font size="2">371</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Lebanon</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">黎巴嫩</font></td>
                      <td align="middle" width="20%"><font size="2">LB</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">961</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Lesotho</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">莱索托</font></td>
                      <td align="middle" width="20%"><font size="2">LS</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">266</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Liberia</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">利比里亚</font></td>
                      <td align="middle" width="20%"><font size="2">LR</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">231</font></td>
                      <td align="middle" width="14%"><font size="2">-8</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Libya</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">利比亚</font></td>
                      <td align="middle" width="20%"><font size="2">LY</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">218</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Liechtenstein</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">列支敦士登</font></td>
                      <td align="middle" width="20%"><font size="2">LI</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">423</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%">Lithuania</td>
                      <td align="middle" width="24%"><font size="2">立陶宛</font></td>
                      <td align="middle" width="20%"><font size="2">LT</font></td>
                      <td align="middle" width="15%"><font size="2">370</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Luxembourg</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">卢森堡</font></td>
                      <td align="middle" width="20%"><font size="2">LU</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">352</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Macao</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">澳门</font></td>
                      <td align="middle" width="20%"><font size="2">MO</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">853</font></td>
                      <td align="middle" width="14%"><font size="2">0</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Madagascar</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">马达加斯加</font></td>
                      <td align="middle" width="20%"><font size="2">MG</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">261</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Malawi</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">马拉维</font></td>
                      <td align="middle" width="20%"><font size="2">MW</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">265</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Malaysia</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">马来西亚</font></td>
                      <td align="middle" width="20%"><font size="2">MY</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">60</font></td>
                      <td align="middle" width="14%"><font size="2">-0.5</font></td></tr>
                    <tr>
                      <td width="27%">Maldives</td>
                      <td align="middle" width="24%"><font color="#000000" size="2">马尔代夫</font></td>
                      <td align="middle" width="20%"><font size="2">MV</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">960</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Mali</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">马里</font></td>
                      <td align="middle" width="20%"><font size="2">ML</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">223</font></td>
                      <td align="middle" width="14%"><font size="2">-8</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Malta</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">马耳他</font></td>
                      <td align="middle" width="20%"><font size="2">MT</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">356</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Mariana Is</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">马里亚那群岛</font></td>
                      <td align="middle" width="20%"><font size="2">&nbsp;</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">1670</font></td>
                      <td align="middle" width="14%"><font size="2">+1</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Martinique</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">马提尼克</font></td>
                      <td align="middle" width="20%"><font size="2">&nbsp;</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">596</font></td>
                      <td align="middle" width="14%"><font size="2">-12</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Mauritius</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">毛里求斯</font></td>
                      <td align="middle" width="20%"><font size="2">MU</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">230</font></td>
                      <td align="middle" width="14%"><font size="2">-4</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Mexico</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">墨西哥</font></td>
                      <td align="middle" width="20%"><font size="2">MX</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">52</font></td>
                      <td align="middle" width="14%"><font size="2">-15</font></td></tr>
                    <tr>
                      <td width="27%">Moldova, Republic of </td>
                      <td align="middle" width="24%"><font size="2">摩尔多瓦</font></td>
                      <td align="middle" width="20%"><font size="2">MD</font></td>
                      <td align="middle" width="15%"><font size="2">373</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Monaco</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">摩纳哥</font></td>
                      <td align="middle" width="20%"><font size="2">MC</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">377</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%">Mongolia </td>
                      <td align="middle" width="24%"><font size="2">蒙古</font></td>
                      <td align="middle" width="20%"><font size="2">MN</font></td>
                      <td align="middle" width="15%"><font size="2">976</font></td>
                      <td align="middle" width="14%"><font size="2">0</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Montserrat 
                      Is</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">蒙特塞拉特岛</font></td>
                      <td align="middle" width="20%"><font size="2">MS</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">1664</font></td>
                      <td align="middle" width="14%"><font size="2">-12</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Morocco</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">摩洛哥</font></td>
                      <td align="middle" width="20%"><font size="2">MA</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">212</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Mozambique</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">莫桑比克</font></td>
                      <td align="middle" width="20%"><font size="2">MZ</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">258</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%">Namibia </td>
                      <td align="middle" width="24%"><font size="2">纳米比亚</font></td>
                      <td align="middle" width="20%"><font size="2">NA</font></td>
                      <td align="middle" width="15%"><font size="2">264</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Nauru</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">瑙鲁</font></td>
                      <td align="middle" width="20%"><font size="2">NR</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">674</font></td>
                      <td align="middle" width="14%"><font size="2">+4</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Nepal</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">尼泊尔</font></td>
                      <td align="middle" width="20%"><font size="2">NP</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">977</font></td>
                      <td align="middle" width="14%"><font size="2">-2.3</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Netheriands 
                        Antilles</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">荷属安的列斯</font></td>
                      <td align="middle" width="20%"><font size="2">&nbsp;</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">599</font></td>
                      <td align="middle" width="14%"><font size="2">-12</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Netherlands</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">荷兰</font></td>
                      <td align="middle" width="20%"><font size="2">NL</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">31</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">New 
Zealand</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">新西兰</font></td>
                      <td align="middle" width="20%"><font size="2">NZ</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">64</font></td>
                      <td align="middle" width="14%"><font size="2">+4</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Nicaragua</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">尼加拉瓜</font></td>
                      <td align="middle" width="20%"><font size="2">NI</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">505</font></td>
                      <td align="middle" width="14%"><font size="2">-14</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Niger</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">尼日尔</font></td>
                      <td align="middle" width="20%"><font face="宋体" size="2">NE</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">227</font></td>
                      <td align="middle" width="14%"><font size="2">-8</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Nigeria</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">尼日利亚</font></td>
                      <td align="middle" width="20%"><font size="2">NG</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">234</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%">North Korea</td>
                      <td align="middle" width="24%"><font size="2">朝鲜</font></td>
                      <td align="middle" width="20%"><font size="2">KP</font></td>
                      <td align="middle" width="15%"><font size="2">850</font></td>
                      <td align="middle" width="14%"><font size="2">+1</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Norway</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">挪威</font></td>
                      <td align="middle" width="20%"><font size="2">NO</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">47</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Oman</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">阿曼</font></td>
                      <td align="middle" width="20%"><font size="2">OM</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">968</font></td>
                      <td align="middle" width="14%"><font size="2">-4</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Pakistan</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">巴基斯坦</font></td>
                      <td align="middle" width="20%"><font size="2">PK</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">92</font></td>
                      <td align="middle" width="14%"><font size="2">-2.3</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Panama</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">巴拿马</font></td>
                      <td align="middle" width="20%"><font size="2">PA</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">507</font></td>
                      <td align="middle" width="14%"><font size="2">-13</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Papua New 
                        Cuinea</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">巴布亚新几内亚</font></td>
                      <td align="middle" width="20%"><font size="2">PG</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">675</font></td>
                      <td align="middle" width="14%"><font size="2">+2</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Paraguay</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">巴拉圭</font></td>
                      <td align="middle" width="20%"><font size="2">PY</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">595</font></td>
                      <td align="middle" width="14%"><font size="2">-12</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Peru</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">秘鲁</font></td>
                      <td align="middle" width="20%"><font size="2">PE</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">51</font></td>
                      <td align="middle" width="14%"><font size="2">-13</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Philippines</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">菲律宾</font></td>
                      <td align="middle" width="20%"><font size="2">PH</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">63</font></td>
                      <td align="middle" width="14%"><font size="2">0</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Poland</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">波兰</font></td>
                      <td align="middle" width="20%"><font size="2">PL</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">48</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">French 
                        Polynesia</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">法属玻利尼西亚</font></td>
                      <td align="middle" width="20%"><font size="2">PF</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">689</font></td>
                      <td align="middle" width="14%"><font size="2">+3</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Portugal</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">葡萄牙</font></td>
                      <td align="middle" width="20%"><font size="2">PT</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">351</font></td>
                      <td align="middle" width="14%"><font size="2">-8</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Puerto 
Rico</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">波多黎各</font></td>
                      <td align="middle" width="20%"><font size="2">PR</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">1787</font></td>
                      <td align="middle" width="14%"><font size="2">-12</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Qatar</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">卡塔尔</font></td>
                      <td align="middle" width="20%"><font size="2">QA</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">974</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Reunion</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">留尼旺</font></td>
                      <td align="middle" width="20%"><font size="2">&nbsp;</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">262</font></td>
                      <td align="middle" width="14%"><font size="2">-4</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Romania</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">罗马尼亚</font></td>
                      <td align="middle" width="20%"><font size="2">RO</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">40</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%">Russia</td>
                      <td align="middle" width="24%"><font size="2">俄罗斯</font></td>
                      <td align="middle" width="20%"><font size="2">RU</font></td>
                      <td align="middle" width="15%"><font size="2">7</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%">Saint Lueia</td>
                      <td align="middle" width="24%"><font size="2">圣卢西亚</font></td>
                      <td align="middle" width="20%"><font size="2">LC</font></td>
                      <td align="middle" width="15%"><font size="2">1758</font></td>
                      <td align="middle" width="14%"><font size="2">-12</font></td></tr>
                    <tr>
                      <td width="27%">Saint Vincent</td>
                      <td align="middle" width="24%"><font size="2">圣文森特岛</font></td>
                      <td align="middle" width="20%"><font size="2">VC</font></td>
                      <td align="middle" width="15%"><font size="2">1784</font></td>
                      <td align="middle" width="14%"><font size="2">-12</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Samoa 
                      Eastern</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">东萨摩亚(美)</font></td>
                      <td align="middle" width="20%"><font size="2">&nbsp;</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">684</font></td>
                      <td align="middle" width="14%"><font size="2">-19</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Samoa 
                      Western</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">西萨摩亚</font></td>
                      <td align="middle" width="20%"><font size="2">&nbsp;</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">685</font></td>
                      <td align="middle" width="14%"><font size="2">-19</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">San Marino</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">圣马力诺</font></td>
                      <td align="middle" width="20%"><font size="2">SM</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">378</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Sao Tome and 
                        Principe</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">圣多美和普林西比</font></td>
                      <td align="middle" width="20%"><font size="2">ST</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">239</font></td>
                      <td align="middle" width="14%"><font size="2">-8</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Saudi 
                      Arabia</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">沙特阿拉伯</font></td>
                      <td align="middle" width="20%"><font size="2">SA</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">966</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Senegal</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">塞内加尔</font></td>
                      <td align="middle" width="20%"><font size="2">SN</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">221</font></td>
                      <td align="middle" width="14%"><font size="2">-8</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Seychelles</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">塞舌尔</font></td>
                      <td align="middle" width="20%"><font size="2">SC</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">248</font></td>
                      <td align="middle" width="14%"><font size="2">-4</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Sierra 
                      Leone</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">塞拉利昂</font></td>
                      <td align="middle" width="20%"><font size="2">SL</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">232</font></td>
                      <td align="middle" width="14%"><font size="2">-8</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Singapore</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">新加坡</font></td>
                      <td align="middle" width="20%"><font size="2">SG</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">65</font></td>
                      <td align="middle" width="14%"><font size="2">+0.3</font></td></tr>
                    <tr>
                      <td width="27%">Slovakia</td>
                      <td align="middle" width="24%"><font color="#000000" size="2">斯洛伐克</font></td>
                      <td align="middle" width="20%"><font size="2">SK</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">421</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%">Slovenia</td>
                      <td align="middle" width="24%"><font size="2">斯洛文尼亚</font></td>
                      <td align="middle" width="20%"><font size="2">SI</font></td>
                      <td align="middle" width="15%"><font size="2">386</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td height="18" width="27%"><font color="#000000">Solomon 
                        Is</font></td>
                      <td align="middle" height="18" width="24%"><font color="#000000" size="2">所罗门群岛</font></td>
                      <td align="middle" height="18" width="20%"><font size="2">SB</font></td>
                      <td align="middle" height="18" width="15%"><font color="#000000" size="2">677</font></td>
                      <td align="middle" height="18" width="14%"><font size="2">+3</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Somali</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">索马里</font></td>
                      <td align="middle" width="20%"><font size="2">SO</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">252</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%">South Africa</td>
                      <td align="middle" width="24%"><font size="2">南非</font></td>
                      <td align="middle" width="20%"><font size="2">ZA</font></td>
                      <td align="middle" width="15%"><font size="2">27</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Spain</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">西班牙</font></td>
                      <td align="middle" width="20%"><font size="2">ES</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">34</font></td>
                      <td align="middle" width="14%"><font size="2">-8</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Sri Lanka</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">斯里兰卡</font></td>
                      <td align="middle" width="20%"><font size="2">LK</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">94</font></td>
                      <td align="middle" width="14%"><font size="2">0</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">St.Lucia</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">圣卢西亚</font></td>
                      <td align="middle" width="20%"><font size="2">LC</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">1758</font></td>
                      <td align="middle" width="14%"><font size="2">-12</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">St.Vincent</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">圣文森特</font></td>
                      <td align="middle" width="20%"><font size="2">VC</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">1784</font></td>
                      <td align="middle" width="14%"><font size="2">-12</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Sudan</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">苏丹</font></td>
                      <td align="middle" width="20%"><font size="2">SD</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">249</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Suriname</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">苏里南</font></td>
                      <td align="middle" width="20%"><font size="2">SR</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">597</font></td>
                      <td align="middle" width="14%"><font size="2">-11.3</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Swaziland</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">斯威士兰</font></td>
                      <td align="middle" width="20%"><font size="2">SZ</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">268</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Sweden</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">瑞典</font></td>
                      <td align="middle" width="20%"><font size="2">SE</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">46</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Switzerland</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">瑞士</font></td>
                      <td align="middle" width="20%"><font size="2">CH</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">41</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Syria</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">叙利亚</font></td>
                      <td align="middle" width="20%"><font size="2">SY</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">963</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%">Taiwan</td>
                      <td align="middle" width="24%"><font size="2">台湾省</font></td>
                      <td align="middle" width="20%"><font size="2">TW</font></td>
                      <td align="middle" width="15%"><font size="2">886</font></td>
                      <td align="middle" width="14%"><font size="2">0</font></td></tr>
                    <tr>
                      <td width="27%">Tajikstan</td>
                      <td align="middle" width="24%"><font size="2">塔吉克斯坦</font></td>
                      <td align="middle" width="20%"><font size="2">TJ</font></td>
                      <td align="middle" width="15%"><font size="2">992</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Tanzania</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">坦桑尼亚</font></td>
                      <td align="middle" width="20%"><font size="2">TZ</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">255</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Thailand</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">泰国</font></td>
                      <td align="middle" width="20%"><font size="2">TH</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">66</font></td>
                      <td align="middle" width="14%"><font size="2">-1</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Togo</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">多哥</font></td>
                      <td align="middle" width="20%"><font size="2">TG</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">228</font></td>
                      <td align="middle" width="14%"><font size="2">-8</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Tonga</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">汤加</font></td>
                      <td align="middle" width="20%"><font size="2">TO</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">676</font></td>
                      <td align="middle" width="14%"><font size="2">+4</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Trinidad and Tobago</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">特立尼达和多巴哥</font></td>
                      <td align="middle" width="20%"><font size="2">TT</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">1809</font></td>
                      <td align="middle" width="14%"><font size="2">-12</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Tunisia</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">突尼斯</font></td>
                      <td align="middle" width="20%"><font size="2">TN</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">216</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Turkey</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">土耳其</font></td>
                      <td align="middle" width="20%"><font size="2">TR</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">90</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%">Turkmenistan</td>
                      <td align="middle" width="24%"><font size="2">土库曼斯坦</font></td>
                      <td align="middle" width="20%"><font size="2">TM</font></td>
                      <td align="middle" width="15%"><font size="2">993</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Uganda</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">乌干达</font></td>
                      <td align="middle" width="20%"><font size="2">UG</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">256</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%">Ukraine</td>
                      <td align="middle" width="24%"><font size="2">乌克兰</font></td>
                      <td align="middle" width="20%"><font size="2">UA</font></td>
                      <td align="middle" width="15%"><font size="2">380</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">United Arab Emirates</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">阿拉伯联合酋长国</font></td>
                      <td align="middle" width="20%"><font size="2">AE</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">971</font></td>
                      <td align="middle" width="14%"><font size="2">-4</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">United Kingdom</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">英国</font></td>
                      <td align="middle" width="20%"><font size="2">GB</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">44</font></td>
                      <td align="middle" width="14%"><font size="2">-8</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">United States</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">美国</font></td>
                      <td align="middle" width="20%"><font size="2">US</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">1</font></td>
                      <td align="middle" width="14%"><font size="2">-13</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Uruguay</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">乌拉圭</font></td>
                      <td align="middle" width="20%"><font size="2">UY</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">598</font></td>
                      <td align="middle" width="14%"><font size="2">-10.3</font></td></tr>
                    <tr>
                      <td width="27%">Uzbekistan</td>
                      <td align="middle" width="24%"><font size="2">乌兹别克斯坦</font></td>
                      <td align="middle" width="20%"><font size="2">UZ</font></td>
                      <td align="middle" width="15%"><font size="2">233</font></td>
                      <td align="middle" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Venezuela</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">委内瑞拉</font></td>
                      <td align="middle" width="20%"><font size="2">VE</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">58</font></td>
                      <td align="middle" width="14%"><font size="2">-12.3</font></td></tr>
                    <tr>
                      <td width="27%">Vietnam</td>
                      <td align="middle" width="24%"><font size="2">越南</font></td>
                      <td align="middle" width="20%"><font size="2">VN</font></td>
                      <td align="middle" width="15%"><font size="2">84</font></td>
                      <td align="middle" width="14%"><font size="2">-1</font></td></tr>
                    <tr>
                      <td height="19" width="27%"><font color="#000000">Yemen</font></td>
                      <td align="middle" height="19" width="24%"><font color="#000000" size="2">也门</font></td>
                      <td align="middle" height="19" width="20%"><font size="2">YE</font></td>
                      <td align="middle" height="19" width="15%"><font color="#000000" size="2">967</font></td>
                      <td align="middle" height="19" width="14%"><font size="2">-5</font></td></tr>
                    <tr>
                      <td width="27%">Yugoslavia</td>
                      <td align="middle" width="24%"><font color="#000000" size="2">南斯拉夫</font></td>
                      <td align="middle" width="20%"><font size="2">YU</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">381</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%">Zimbabwe</td>
                      <td align="middle" width="24%"><font color="#000000" size="2">津巴布韦</font></td>
                      <td align="middle" width="20%"><font size="2">ZW</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">263</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td></tr>
                    <tr>
                      <td width="27%">Zaire</td>
                      <td align="middle" width="24%"><font size="2">扎伊尔</font></td>
                      <td align="middle" width="20%"><font size="2">ZR</font></td>
                      <td align="middle" width="15%"><font size="2">243</font></td>
                      <td align="middle" width="14%"><font size="2">-7</font></td></tr>
                    <tr>
                      <td width="27%"><font color="#000000">Zambia</font></td>
                      <td align="middle" width="24%"><font color="#000000" size="2">赞比亚</font></td>
                      <td align="middle" width="20%"><font size="2">ZM</font></td>
                      <td align="middle" width="15%"><font color="#000000" size="2">260</font></td>
                      <td align="middle" width="14%"><font size="2">-6</font></td>
                     </tr>`

	ss := strings.Split(html, `<tr>`)
	countries := make([]*Country, 0)
	for _, s1 := range ss {
		if strings.Contains(s1, `<td width="27%">`) {
			//fmt.Println(strings.Split(s1,`<td width="27%">`))
			area := strings.ReplaceAll(strings.Split(strings.Split(s1, `<td width="27%">`)[1], `</td>`)[0], `<font color="#000000">`, "")
			area = strings.ReplaceAll(area, `</font>`, "")
			area = strings.ReplaceAll(area, `
`, " ")
			area = strings.ReplaceAll(strings.ReplaceAll(area, `                        `, " "), `   `, "")
			//fmt.Println(area)
			areaName := strings.ReplaceAll(strings.Split(strings.Split(s1, `<td align="middle" width="24%">`)[1], `</td>`)[0], `<font color="#000000" size="2">`, "")
			areaName = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(areaName, `</font>`, ""), `<fontsize="2">`, ""), ` `, "")
			areaName = strings.ReplaceAll(areaName, `
`, " ")
			areaName = strings.ReplaceAll(areaName, `<fontsize="2">`, "")
			//fmt.Println(areaName) United Kingdom

			shortName := strings.ReplaceAll(strings.Split(strings.Split(s1, `<td align="middle" width="20%">`)[1], `</td>`)[0], `<font size="2">`, "")
			shortName = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(shortName, `</font>`, ""), `<fontsize="2">`, ""), ` `, "")
			shortName = strings.ReplaceAll(shortName, `
`, " ")
			shortName = strings.ReplaceAll(shortName, `<fontsize="2">`, "")

			//                     United Kiongdom
			code := strings.ReplaceAll(strings.Split(strings.Split(s1, `<td align="middle" width="15%">`)[1], `</td>`)[0], `<font color="#000000" size="2">`, "")
			code = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(code, `</font>`, ""), " ", ""), `<fontsize="2">`, "")
			code = strings.ReplaceAll(code, `
`, " ")
			//fmt.Println(code)
			countries = append(countries, &Country{Name: areaName, EnglishName: area, Code: code, ShortName: shortName})
		}
	}
	m := make(map[string]string)
	_ = m
	db := orm.NewMySQLDB("cti_link_conf_db", "ctiLinkMysql", "ctiLinkMysql", "localhost", 3306)
	db.Open()
	defer db.Close()
	pkg.ReadLine("sql.csv", func(s string) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err, s)
			}
		}()
		ss := strings.Split(s, ",")
		exit := Query(ss[0], db)
		if !exit {
			pkg.ToTxt(s, "No.txt")
		}
	})

	/*countries1 := make([]*Country,0)
	for _,c := range countries{
		exit := false
		if _, ok := m[c.EnglishName]; ok {
			exit = true
		}
		if _, ok := m[c.EnglishName+" "]; ok{
			exit = true
		}
		if !exit{
			//fmt.Println(c)
			countries1 = append(countries1,c)
		}
	}
	Detail(countries1)
	*/
	/*for _,c:=range countries1{
		fmt.Println(c)
	}*/
	/*阿森松,Ascension,247
	开曼群岛,Cayman,1345
	法属圭亚那,French Guiana,594
	科特迪瓦,Ivory  Coast,225
	马里亚那群岛,Mariana Is,1670
	马提尼克,Martinique,596
	荷属安的列斯,NetheriandsAntilles,599
	留尼旺,Reunion,262
	东萨摩亚(美),Samoa Eastern,684
	西萨摩亚,Samoa Western,685
	南斯拉夫,Yugoslavia,381
	扎伊尔,Zaire,243
	*/
	/*for k,_ := range m{
		pkg.ToTxt(k,"name.txt")
	}*/
}

func Query(prefix string, db orm.DB) bool {
	var rows *sql.Rows
	var err error
	if strings.Index(prefix, "1") == 0 {
		rows, err = db.DB().Query("select * from cti_link_international_area_code where prefix ='" + prefix + "'")
	} else {
		rows, err = db.DB().Query("select * from cti_link_international_area_code where country_code ='" + prefix + "'")
	}

	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	areaCodes := make([]AreaCode, 0)

	for rows.Next() {
		var areaCode AreaCode
		if err := rows.Scan(&areaCode.Id, &areaCode.Prefix, &areaCode.AreaCode, &areaCode.City, &areaCode.CountryCode, &areaCode.Country); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		fmt.Println(areaCode)
		areaCodes = append(areaCodes, areaCode)
	}
	rows.Close()
	if len(areaCodes) > 0 {
		return true
	} else {
		return false
	}
}

type AreaCode struct {
	Id          int64
	Prefix      string
	AreaCode    string
	City        string
	CountryCode string
	Country     string
}

type Country struct {
	Name        string
	EnglishName string
	ShortName   string
	Code        string
	Areas       map[string]string
}

func Detail(countries []*Country) {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("headless", false),
		chromedp.WindowSize(1280, 1024),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	var res string
	tasks := chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("step1")
			for _, country := range countries {
				chromedp.Navigate("https://countrycode.org/" + country.ShortName).Do(ctx)
				chromedp.OuterHTML(`body`, &res, chromedp.ByQuery).Do(ctx)
				if !strings.Contains(res, `We are sorry. The page you requested cannot be found`) {
					chromedp.OuterHTML(`collapseCityCodes`, &res, chromedp.ByID).Do(ctx)
					country.Areas = GetCode(res)
					//fmt.Println(country.Name)
					time.Sleep(2 * time.Second)
					for k, v := range country.Areas {
						pkg.ToTxt(country.Name+","+country.EnglishName+","+country.Code+","+k+","+v, "codes.csv")
					}
				} else {
					fmt.Println(country.Name + "," + country.EnglishName + "," + country.Code)
				}

			}
			return nil
		}),
	}
	fmt.Println("start run chromedp tasks!")
	err := chromedp.Run(ctx,
		tasks,
	)
	if err != nil {
		log.Println(err)
	}
}

func GetCode(s string) map[string]string {
	/*s := `<div id="collapseCityCodes" class="table-responsive text-uppercase collapse in">
	              <div class="bootstrap-table"><div class="fixed-table-toolbar"></div><div class="fixed-table-container"><div class="fixed-table-header" style="display: block;"><table></table></div><div class="fixed-table-body"><div class="fixed-table-loading" style="top: 37px;">Loading, please wait...</div><table data-toggle="table" data-classes="table table-hover table-striped" data-sort-name="citycode" data-sort-order="desc" class="table table-hover table-striped">
	                  <thead>
	                      <tr><th class="cities-col" style=""><div class="th-inner sortable">
	                              <i class="fa fa-globe i_blue"></i> City
	                              <a href="/countryCode/downloadCityCodes?country=FR" class="pull-right"><i class="fa fa-download"></i></a>
	                      </div><div class="fht-cell"></div></th><th style=""><div class="th-inner sortable">Dial Codes</div><div class="fht-cell"></div></th></tr>
	              </thead>
	              <tbody><tr data-index="0"><td class="cities-col" style="">Northeast France</td><td style="">+33-3</td></tr><tr data-index="1"><td class="cities-col" style="">Northwest France</td><td style="">+33-2</td></tr><tr data-index="2"><td class="cities-col" style="">Paris and Paris Region</td><td style="">+33-1</td></tr><tr data-index="3"><td class="cities-col" style="">Southeast France and Corsica</td><td style="">+33-4</td></tr><tr data-index="4"><td class="cities-col" style="">Southwest France</td><td style="">+33-5</td></tr></tbody>
	              </table></div><div class="fixed-table-pagination" style="display: none;"></div></div></div><div class="clearfix"></div>
	          </div>
	`*/
	tbody := strings.Split(s, `<tbody>`)[1]
	trs := strings.Split(tbody, `<tr`)
	m := make(map[string]string)
	for _, tr := range trs {
		tr = strings.Split(tr, `</tr>`)[0]
		if strings.Contains(tr, `data-index`) {
			ss := strings.Split(tr, `style="">`)
			name := strings.ReplaceAll(ss[1], `</td>`, "")
			name = strings.ReplaceAll(name, `<td `, "")
			code := strings.ReplaceAll(ss[2], `</td>`, "")
			//fmt.Println(name,"=",code)
			m[name] = code
		}
	}
	return m
}
