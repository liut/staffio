/*
 * jQuery pretty date plug-in 1.0.0
 *
 * http://bassistance.de/jquery-plugins/jquery-plugin-prettydate/
 *
 * Based on John Resig's prettyDate http://ejohn.org/blog/javascript-pretty-date
 *
 * Copyright (c) 2009 Jörn Zaefferer
 *
 * $Id: 3b1e80b0c4f206062103fe3bace832449e473727 $
 *
 * Dual licensed under the MIT and GPL licenses:
 *   http://www.opensource.org/licenses/mit-license.php
 *   http://www.gnu.org/licenses/gpl.html
 */

(function($) {

$.prettyDate = {

	template: function(source, params) {
		if ( arguments.length == 1 )
			return function() {
				var args = $.makeArray(arguments);
				args.unshift(source);
				return $.prettyDate.template.apply( this, args );
			};
		if ( arguments.length > 2 && params.constructor != Array  ) {
			params = $.makeArray(arguments).slice(1);
		}
		if ( params.constructor != Array ) {
			params = [ params ];
		}
		$.each(params, function(i, n) {
			source = source.replace(new RegExp("\\{" + i + "\\}", "g"), n);
		});
		return source;
	},

	now: function() {
		return new Date();
	},

	// Takes an ISO time and returns a string representing how
	// long ago the date represents.
	format: function(time) {
		if (typeof time === "undefined") return;
		var date, re = /(\d{4})[\/-](\d{2})[\/-](\d{2})[\sT](\d{2})\:(\d{2})\:(\d{2})(\.(\d{1,3}))?/,
		r = time.match(re), ms_time, diff, day_diff, now = $.prettyDate.now(),
		serverTZ = window.JS_TIMEZONE_OFFSET ? window.JS_TIMEZONE_OFFSET : now.getTimezoneOffset(), offset = (serverTZ - now.getTimezoneOffset()) * 60 * 1000;

		//console.log('serverTZ', serverTZ, 'offset', offset, 'r:', r);

		if (r !== null) {
			ms_time = Date.UTC(r[1],r[2]-1,r[3],r[4],r[5],r[6], r[8] ? r[8] : 0) + serverTZ * 60 * 1000;
		} else {
			date = new Date((time || "").replace(/-/g,"/").replace(/[TZ]/g," "));
			ms_time = date.getTime();
		}
		diff = (now.getTime() - ms_time /*+ offset*/) / 1000,
		day_diff = Math.floor(diff / 86400);

		//console.log(now.getTime(),ms_time, diff, day_diff, now, new Date(ms_time + offset));

		var messages = $.prettyDate.messages;
		if ( isNaN(day_diff) || day_diff < 0 )
			return time;
		if ( day_diff >= 31 ) {
			date = new Date(ms_time);
			var year = date.getFullYear();
			if (now.getFullYear() == year) {
				return (date.getMonth() + 1) + messages.month + date.getDate() + messages.day;
			}
			return year + "-" +(date.getMonth() + 1) + "-" + date.getDate();
			//return date.toLocaleDateString();
		}

		return day_diff == 0 && (
				diff < 60 && messages.now ||
				diff < 120 && messages.minute ||
				diff < 3600 && messages.minutes(Math.floor( diff / 60 )) ||
				diff < 7200 && messages.hour ||
				diff < 86400 && messages.hours(Math.floor( diff / 3600 ))) ||
			day_diff == 1 && messages.yesterday ||
			day_diff == 2 && messages.the_day_before_yesterday ||
			day_diff < 7 && messages.days(day_diff) ||
			day_diff < 31 && messages.weeks(Math.ceil( day_diff / 7 ));
	}
	,
	count: 0
	,
	interval: null

};

$.prettyDate.messages = {
	month: "月",
	day: "日",
	now: "刚才",
	minute: "1分钟前",
	minutes: $.prettyDate.template("{0}分钟前"),
	hour: "1小时前",
	hours: $.prettyDate.template("{0}小时前"),
	yesterday: "昨天",
	the_day_before_yesterday: "前天",
	days: $.prettyDate.template("{0}天前"),
	weeks: $.prettyDate.template("{0}周前")
};

$.fn.prettyDate = function(options) {
	options = $.extend({
		value: function() {
			return $(this).attr("title");
		},
		interval: 10000
	}, options);
	var elements = this;
	function format() {
		elements.each(function() {
			var date = $.prettyDate.format(options.value.apply(this));
			if ( date && $(this).text() != date )
				$(this).text( date );
		});
		//$.prettyDate.count ++;
		//if ($.prettyDate.count > 2) {clearInterval($.prettyDate.interval);};
	}
	format();
	if (options.interval)
		/*$.prettyDate.interval =*/ setInterval(format, options.interval);
	return this;
};

if(typeof $.fn.fmatter !== "undefined") {
	$.fn.fmatter.prettydate = function (cellval, opts, act) {
		var date = $.prettyDate.format(cellval);
		if ( date ) return "<span title=\""+ cellval +"\" >"+ date + "</span>";
		return $.fn.fmatter.defaultFormat(cellval, opts);
	};
}

})(jQuery);
