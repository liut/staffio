var jQuery = require('jquery')

global.$ = jQuery
global.jQuery = jQuery

require('bootstrap')
// require('bootstrap-select')
// require('bootstrap-validator')
require('bootstrapValidator/dist/js/bootstrapValidator')
require('bootstrapValidator/dist/js/language/zh_CN')
require('x-editable/dist/bootstrap3-editable/js/bootstrap-editable')

require('./plugin/pwstrength')
var Password = require('./plugin/password')
global.Password = Password
require('./plugin/jquery.prettydate')
// require('moment/locale/zh-cn');

global.Dust = require('./plugin/dust')
require( 'datatables.net' )( global, jQuery );
