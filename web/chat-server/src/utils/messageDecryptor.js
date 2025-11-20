/**
 * 消息解密辅助工具
 * 用于在获取消息列表后批量解密
 */

import { acceptSession, receiveAndDecryptMessage, hasSession } from '@/crypto';
import store from '../store';

/**
 * 解密单条消息
 * @param {Object} message - 消息对象
 * @returns {Promise<Object>} 解密后的消息对象
 */
export async function decryptMessage(message) {
  // 如果不是加密消息，直接返回
  if (!message.is_encrypted) {
    return message;
  }

  // 如果没有主密钥，无法解密
  if (!store.state.masterKey) {
    console.warn('未启用加密，无法解密消息');
    return {
      ...message,
      content: '[加密消息 - 需要启用加密功能]',
    };
  }

  try {
    // 确定发送者（对方）ID
    const contactId =
      message.send_id === store.state.userInfo.uuid
        ? message.receive_id
        : message.send_id;
    
    // 获取消息 ID（可能是 uuid、message_id 或其他字段名）
    // 注意：后端可能返回带尾随空格的 UUID，需要 trim
    let messageId = message.uuid || message.message_id || message.id;
    if (messageId) {
      messageId = messageId.trim(); // 去除首尾空格
    }
    
    // 如果消息是自己发送的，尝试从 IndexedDB 中读取明文（发送方在发送时保存的）
    if (message.send_id === store.state.userInfo.uuid) {
      console.log('ℹ️ 消息是自己发送的，尝试从 IndexedDB 读取明文');
      
      if (!messageId) {
        console.warn('⚠️ 消息 ID 不存在，无法从 IndexedDB 读取明文');
        return {
          ...message,
          content: message.content || '[已发送]',
          _is_sent: true,
        };
      }
      
      try {
        const { get, STORES, getAll } = await import('@/crypto/cryptoStore');
        console.log('🔍 尝试从 IndexedDB 读取明文，messageId:', messageId, 'length:', messageId.length);
        
        // 先尝试直接获取（使用 trim 后的 key）
        let sentMessage = await get(STORES.SENT_MESSAGES, messageId);
        
        // 如果找不到，尝试使用原始 key（可能之前保存时没有 trim）
        if (!sentMessage && message.uuid) {
          const originalKey = message.uuid;
          if (originalKey !== messageId) {
            console.log('⚠️ 使用 trim 后的 key 找不到，尝试使用原始 key:', originalKey, 'length:', originalKey.length);
            sentMessage = await get(STORES.SENT_MESSAGES, originalKey);
          }
        }
        
        // 如果还是找不到，尝试获取所有记录来调试
        if (!sentMessage) {
          console.log('⚠️ 直接获取失败，尝试获取所有记录来调试...');
          const allSentMessages = await getAll(STORES.SENT_MESSAGES);
          console.log('📋 IndexedDB 中所有发送方明文记录:', {
            total: allSentMessages.length,
            message_ids: allSentMessages.map(m => ({
              id: m.message_id,
              length: m.message_id?.length,
              has_spaces: m.message_id?.includes(' '),
            })),
            searching_for: messageId,
            searching_length: messageId.length,
          });
          
          // 尝试在所有记录中查找匹配的（忽略空格）
          const matched = allSentMessages.find(m => 
            m.message_id && m.message_id.trim() === messageId
          );
          if (matched) {
            console.log('✅ 通过模糊匹配找到明文（忽略空格）');
            sentMessage = matched;
          }
        }
        
        if (sentMessage && sentMessage.plaintext) {
          console.log('✅ 从 IndexedDB 读取到发送方的明文:', {
            message_id: sentMessage.message_id,
            plaintext_preview: sentMessage.plaintext.substring(0, 20),
          });
          return {
            ...message,
            content: sentMessage.plaintext, // 使用保存的明文
            _is_sent: true,
          };
        } else {
          console.warn('⚠️ 未找到发送方的明文（messageId:', messageId, '），返回原始内容');
          return {
            ...message,
            content: message.content || '[已发送]', // 如果找不到明文，返回原始内容（可能是密文）
            _is_sent: true,
          };
        }
      } catch (error) {
        console.warn('读取发送方明文失败:', error, 'messageId:', messageId);
        return {
          ...message,
          content: message.content || '[已发送]',
          _is_sent: true,
        };
      }
    }

    // 如果是接收的消息（不是自己发送的），先尝试从 IndexedDB 读取明文
    if (message.send_id !== store.state.userInfo.uuid && messageId) {
      console.log('ℹ️ 消息是接收的，尝试从 IndexedDB 读取明文');
      try {
        const { get, STORES, getAll } = await import('@/crypto/cryptoStore');
        console.log('🔍 尝试从 IndexedDB 读取接收方明文，messageId:', messageId, 'length:', messageId.length);
        
        // 先尝试直接获取（使用 trim 后的 key）
        let receivedMessage = await get(STORES.RECEIVED_MESSAGES, messageId);
        
        // 如果找不到，尝试使用原始 key（可能之前保存时没有 trim）
        if (!receivedMessage && message.uuid) {
          const originalKey = message.uuid;
          if (originalKey !== messageId) {
            console.log('⚠️ 使用 trim 后的 key 找不到，尝试使用原始 key:', originalKey, 'length:', originalKey.length);
            receivedMessage = await get(STORES.RECEIVED_MESSAGES, originalKey);
          }
        }
        
        // 如果还是找不到，尝试在所有记录中模糊匹配（忽略空格）
        if (!receivedMessage) {
          const allReceivedMessages = await getAll(STORES.RECEIVED_MESSAGES);
          const matched = allReceivedMessages.find(m => 
            m.message_id && m.message_id.trim() === messageId
          );
          if (matched) {
            console.log('✅ 通过模糊匹配找到接收方明文（忽略空格）');
            receivedMessage = matched;
          }
        }
        
        if (receivedMessage && receivedMessage.plaintext) {
          console.log('✅ 从 IndexedDB 读取到接收方的明文:', {
            message_id: receivedMessage.message_id,
            plaintext_preview: receivedMessage.plaintext.substring(0, 20),
          });
          return {
            ...message,
            content: receivedMessage.plaintext, // 使用保存的明文
            _is_received: true,
          };
        } else {
          console.log('ℹ️ 未找到接收方的明文，将尝试解密');
        }
      } catch (error) {
        console.warn('读取接收方明文失败:', error);
        // 继续尝试解密
      }
    }

    // 检查会话是否存在
    let sessionExists = await hasSession(contactId);

    // 如果是 PreKeyMessage 且会话不存在，需要先接受会话
    if (message.message_type === 'PreKeyMessage' && !sessionExists) {
      console.log('🔒 接收 PreKeyMessage，建立会话...');
      console.log('PreKeyMessage 数据:', {
        sender_identity_key: message.sender_identity_key,
        sender_ephemeral_key: message.sender_ephemeral_key,
        ratchet_key: message.ratchet_key,
        used_one_time_pre_key_id: message.used_one_time_pre_key_id,
      });
      
      // 检查必需字段
      if (!message.sender_identity_key || !message.sender_ephemeral_key || !message.ratchet_key) {
        throw new Error('PreKeyMessage 缺少必需字段');
      }
      
      await acceptSession(store.state.masterKey, contactId, {
        identity_key: message.sender_identity_key,
        ephemeral_key: message.sender_ephemeral_key,
        ratchet_key: message.ratchet_key,
        used_one_time_pre_key_id: message.used_one_time_pre_key_id,
      });
      
      // 再次检查会话是否存在（确保保存成功）
      sessionExists = await hasSession(contactId);
      if (!sessionExists) {
        throw new Error('会话建立失败：会话未保存');
      }
      console.log('✅ 会话已建立并验证');
    }

    // 解密消息
    // 对于 PreKeyMessage，counter 可能为 null，默认为 0
    const counter = message.counter !== null && message.counter !== undefined ? message.counter : 0;
    const prevCounter = message.prev_counter !== null && message.prev_counter !== undefined ? message.prev_counter : 0;
    
    console.log('📥 接收到的消息数据:', {
      is_encrypted: message.is_encrypted,
      message_type: message.message_type,
      content_length: message.content?.length,
      content_preview: message.content?.substring(0, 20),
      iv_length: message.iv?.length,
      auth_tag_length: message.auth_tag?.length,
      ratchet_key_length: message.ratchet_key?.length,
      counter: counter,
    });
    
    const encryptedData = {
      ratchet_key: message.ratchet_key,
      counter: counter,
      prev_counter: prevCounter,
      ciphertext: message.content, // content 字段存储的是密文
      iv: message.iv,
      auth_tag: message.auth_tag,
    };

    const plaintext = await receiveAndDecryptMessage(contactId, encryptedData);
    console.log('✅ 消息已解密');

    // 保存接收方的明文到 IndexedDB（仅接收方可见，用于历史记录）
    // messageId 已经在上面声明过了，直接使用
    if (messageId && store.state.masterKey) {
      try {
        const { put, STORES } = await import('@/crypto/cryptoStore');
        const receivedMessageData = {
          message_id: messageId,
          plaintext: plaintext,
          contact_id: contactId,
          created_at: Date.now(),
        };
        await put(STORES.RECEIVED_MESSAGES, receivedMessageData);
        console.log('📝 已保存接收方的明文到 IndexedDB:', {
          message_id: messageId,
          plaintext_length: plaintext.length,
          plaintext_preview: plaintext.substring(0, 20),
        });
      } catch (error) {
        console.warn('保存接收方明文失败:', error);
        // 不影响消息解密，只是历史记录可能无法显示
      }
    }

    // 返回解密后的消息
    return {
      ...message,
      content: plaintext, // 替换为明文
      _decrypted: true,
    };
  } catch (error) {
    console.error('解密消息失败:', error);
    return {
      ...message,
      content: '[解密失败: ' + error.message + ']',
      _decrypt_error: true,
    };
  }
}

/**
 * 批量解密消息列表
 * @param {Array} messages - 消息列表
 * @returns {Promise<Array>} 解密后的消息列表
 */
export async function decryptMessageList(messages) {
  if (!messages || messages.length === 0) {
    return messages;
  }

  console.log(`准备解密 ${messages.length} 条消息...`);

  // 如果没有主密钥，无法解密
  if (!store.state.masterKey) {
    console.warn('未启用加密，无法解密消息');
    return messages.map(msg => ({
      ...msg,
      content: msg.is_encrypted ? '[加密消息 - 需要启用加密功能]' : msg.content,
    }));
  }

  // 第一步：先找到所有 PreKeyMessage 并建立会话
  const contactSessions = new Set(); // 记录已建立会话的联系人
  for (const message of messages) {
    if (!message.is_encrypted || message.message_type !== 'PreKeyMessage') {
      continue;
    }

    const contactId =
      message.send_id === store.state.userInfo.uuid
        ? message.receive_id
        : message.send_id;

    // 如果已经为这个联系人建立过会话，跳过
    if (contactSessions.has(contactId)) {
      continue;
    }

    // 检查会话是否已存在
    const sessionExists = await hasSession(contactId);
    if (sessionExists) {
      contactSessions.add(contactId);
      continue;
    }

    // 建立会话
    console.log(`🔒 发现 PreKeyMessage，为 ${contactId} 建立会话...`);
    console.log('PreKeyMessage 数据:', {
      sender_identity_key: message.sender_identity_key,
      sender_ephemeral_key: message.sender_ephemeral_key,
      ratchet_key: message.ratchet_key,
      used_one_time_pre_key_id: message.used_one_time_pre_key_id,
    });

    // 检查必需字段
    if (!message.sender_identity_key || !message.sender_ephemeral_key || !message.ratchet_key) {
      console.error('PreKeyMessage 缺少必需字段，跳过建立会话');
      continue;
    }

    try {
      await acceptSession(store.state.masterKey, contactId, {
        identity_key: message.sender_identity_key,
        identity_key_curve25519: message.sender_identity_key_curve25519, // Curve25519 格式的身份公钥（用于 ECDH）
        ephemeral_key: message.sender_ephemeral_key,
        ratchet_key: message.ratchet_key,
        used_one_time_pre_key_id: message.used_one_time_pre_key_id,
      });

      // 验证会话是否成功建立
      const verified = await hasSession(contactId);
      if (verified) {
        contactSessions.add(contactId);
        console.log(`✅ 会话已建立并验证: ${contactId}`);
      } else {
        console.error(`❌ 会话建立失败: ${contactId}`);
      }
    } catch (error) {
      console.error(`建立会话失败 (${contactId}):`, error);
    }
  }

  // 第二步：解密所有消息
  const decryptedMessages = [];
  for (const message of messages) {
    const decrypted = await decryptMessage(message);
    decryptedMessages.push(decrypted);
  }

  console.log(`✅ ${decryptedMessages.length} 条消息解密完成`);
  return decryptedMessages;
}

