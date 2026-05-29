// WebSocket客户端
(function(global) {
    'use strict';

    let ws = null;
    let reconnectAttempts = 0;
    const maxReconnectAttempts = 5;
    const reconnectDelay = 3000;

    // 初始化WebSocket连接
    global.initWebSocket = function(userId) {
        if (!userId || userId <= 0) {
            console.error('无效的用户ID');
            return;
        }

        // 如果已有连接，先关闭
        if (ws && ws.readyState === WebSocket.OPEN) {
            ws.close();
        }

        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//127.0.0.1:8888/ws?user_id=${userId}`;

        console.log('正在连接WebSocket:', wsUrl);
        ws = new WebSocket(wsUrl);

        ws.onopen = function() {
            console.log('WebSocket连接已建立');
            reconnectAttempts = 0;
            global.showToast('已连接到服务器', 'success');
        };

        ws.onmessage = function(event) {
            try {
                const message = JSON.parse(event.data);
                console.log('收到WebSocket消息:', message);
                handleMessage(message);
            } catch (e) {
                console.error('解析WebSocket消息失败:', e);
            }
        };

        ws.onclose = function(event) {
            console.log('WebSocket连接已关闭:', event.code, event.reason);
            attemptReconnect(userId);
        };

        ws.onerror = function(error) {
            console.error('WebSocket错误:', error);
            global.showToast('连接错误', 'error');
        };
    };

    // 处理收到的消息
    function handleMessage(message) {
        switch (message.type) {
            case 'payment_reminder':
                handlePaymentReminder(message.data);
                break;
            case 'heartbeat':
                // 心跳消息，忽略
                break;
            default:
                console.log('未知消息类型:', message.type);
        }
    }

    // 处理支付提醒
    function handlePaymentReminder(data) {
        const { order_id, start_address, end_address, actual_price, remind_count } = data;

        // 显示提醒弹窗
        showPaymentReminderModal(order_id, start_address, end_address, actual_price, remind_count);
    }

    // 显示支付提醒弹窗
    function showPaymentReminderModal(orderId, startAddress, endAddress, actualPrice, remindCount) {
        // 检查是否已存在弹窗，避免重复显示
        let modal = document.getElementById('payment-reminder-modal');
        if (!modal) {
            // 创建弹窗HTML
            modal = document.createElement('div');
            modal.id = 'payment-reminder-modal';
            modal.innerHTML = `
                <div style="
                    position: fixed;
                    top: 0;
                    left: 0;
                    width: 100%;
                    height: 100%;
                    background: rgba(0,0,0,0.5);
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    z-index: 9999;
                ">
                    <div style="
                        background: white;
                        padding: 30px;
                        border-radius: 15px;
                        width: 90%;
                        max-width: 400px;
                        text-align: center;
                        box-shadow: 0 10px 30px rgba(0,0,0,0.3);
                    ">
                        <div style="
                            font-size: 50px;
                            margin-bottom: 20px;
                        ">⚠️</div>
                        <h3 style="margin: 0 0 15px 0; color: #333;">
                            请及时完成支付
                        </h3>
                        <div id="reminder-content" style="margin-bottom: 20px; color: #666;">
                        </div>
                        <div style="display: flex; gap: 10px; justify-content: center;">
                            <button id="remind-later-btn" style="
                                padding: 10px 25px;
                                border: none;
                                border-radius: 20px;
                                cursor: pointer;
                                font-size: 14px;
                                background: #f0f0f0;
                                color: #666;
                            ">稍后</button>
                            <button id="go-pay-btn" style="
                                padding: 10px 25px;
                                border: none;
                                border-radius: 20px;
                                cursor: pointer;
                                font-size: 14px;
                                background: linear-gradient(135deg, #ff6b00 0%, #ffa500 100%);
                                color: white;
                            ">去支付</button>
                        </div>
                    </div>
                </div>
            `;
            document.body.appendChild(modal);

            // 绑定事件
            document.getElementById('remind-later-btn').onclick = function() {
                modal.style.display = 'none';
            };

            document.getElementById('go-pay-btn').onclick = function() {
                // 跳转到订单页面
                window.location.href = './order.html';
            };
        }

        // 更新内容
        document.getElementById('reminder-content').innerHTML = `
            <p><strong>订单号:</strong> ${orderId}</p>
            <p><strong>行程:</strong> ${startAddress.slice(0, 15)}${startAddress.length > 15 ? '...' : ''} → ${endAddress.slice(0, 15)}${endAddress.length > 15 ? '...' : ''}</p>
            <p><strong>待支付:</strong> ¥${actualPrice.toFixed(2)}</p>
            <p style="color: #ff6b00; font-size: 12px;">（第${remindCount}次提醒）</p>
        `;

        // 显示弹窗
        modal.style.display = 'flex';
    }

    // 尝试重连
    function attemptReconnect(userId) {
        if (reconnectAttempts < maxReconnectAttempts) {
            reconnectAttempts++;
            console.log(`尝试重连 (${reconnectAttempts}/${maxReconnectAttempts})...`);
            setTimeout(function() {
                global.initWebSocket(userId);
            }, reconnectDelay);
        } else {
            console.error('WebSocket重连失败');
            global.showToast('连接断开，请刷新页面', 'error');
        }
    }

    // 关闭WebSocket连接
    global.closeWebSocket = function() {
        if (ws) {
            ws.close();
            ws = null;
            console.log('WebSocket连接已关闭');
        }
    };

})(window);
