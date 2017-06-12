-- init

INSERT INTO oauth_scope(name,label,description,is_default) VALUES('basic', 'Basic', 'Read your Uid (login name) and Nickname', true);
INSERT INTO oauth_scope(name,label,description) VALUES('profile', 'Personal Information', 'Read your GivenName, Surname, Email, etc.');

-- INSERT INTO oauth_client VALUES(1, '1234', 'Demo', 'aabbccdd', 'http://localhost:3000/appauth', '{}', now());
-- INSERT INTO links(title, url, author) VALUES('周报', 'https://weekly.lcgc.work', 'liutao');

/*
INSERT INTO article(title, content, author)
VALUES('`Header` is header', '已在内部registry并经过测试的docker images:

 - `lcgc/python:2.7.12-r0`
 - `lcgc/mariadb:10.1.19`
 - `lcgc/golang:1.7.3-r0`
 - `lcgc/sso-db:v0.13.2`', 'liutao');

INSERT INTO article(title, content, author)
VALUES('周四分享预告：一小时学会使用`Charles`抓网络数据包', '本周分享预告来鸟

上周的热修复是不是只听懂了一小半？没关系，这次是扫盲式培训，手把手教你1小时当场学会抓网络数据包。

效果： 学完之后可以了解到任何一个app是怎么发送数据和接收数据，以及什么样的数据，然后才能展示在界面上。

目的： 帮助经常跟技术打交道的同学沟通更顺畅，当然也可以尝试查看感兴趣的app的接口，做点黑科技，比如投票刷票等等。
适合人群：

1. 完全没有技术基础，扫盲式培训；
2. 不知道工程师们经常讨论的接口是什么；
3. 不了解`Charles`如何来使用的，或者只会使用基本功能的；
4. 前端工程师，app工程师，测试工程师，产品经理，以及其他一切经常与程序员打交道的岗位都应该学习一下；
5. 机会难得，仅此一次；

分享人：郭亚伦<br>
分享地点：陆家嘴会议室<br>
时间：22号周四下午4:00', 'shenshanshan');

-- drop table oauth_client, oauth_access_token, oauth_authorization_code, oauth_scope, oauth_client_user_authorized, articles, links, password_reset, user_log cascade;
-- drop table cas_ticket, http_sessions;
*/
