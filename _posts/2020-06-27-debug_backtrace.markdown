---
title: debug_backtrace
layout: post
category: php
author: 夏泽民
---
Zend/zend_buildin_functions.c +2600
ZEND_FUNCTION(debug_backtrace)

ZEND_API void zend_fetch_debug_backtrace(zval *return_value, int skip_last, int options, int limit) 
<!-- more -->
if (!(ptr = EG(current_execute_data))) {
}

if (!ptr->func || !ZEND_USER_CODE(ptr->func->common.type)) {
		call = ptr;
		ptr = ptr->prev_execute_data;
}

			if (ptr->func && ZEND_USER_CODE(ptr->func->common.type) && (ptr->opline->opcode == ZEND_NEW)) {
				call = ptr;
				ptr = ptr->prev_execute_data;
			}
			
		array_init(&stack_frame);

		ptr = zend_generator_check_placeholder_frame(ptr);
		
				if ((!skip->func || !ZEND_USER_CODE(skip->func->common.type)) &&
		    skip->prev_execute_data &&
		    skip->prev_execute_data->func &&
		    ZEND_USER_CODE(skip->prev_execute_data->func->common.type) &&
		    skip->prev_execute_data->opline->opcode != ZEND_DO_FCALL &&
		    skip->prev_execute_data->opline->opcode != ZEND_DO_ICALL &&
		    skip->prev_execute_data->opline->opcode != ZEND_DO_UCALL &&
		    skip->prev_execute_data->opline->opcode != ZEND_DO_FCALL_BY_NAME &&
		    skip->prev_execute_data->opline->opcode != ZEND_INCLUDE_OR_EVAL) {
			skip = skip->prev_execute_data;
		}