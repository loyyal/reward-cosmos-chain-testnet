import { SigningStargateClient, calculateFee, GasPrice, defaultRegistryTypes } from "@cosmjs/stargate";
import { DirectSecp256k1HdWallet, Registry } from "@cosmjs/proto-signing";
import { stringToPath } from "@cosmjs/crypto";
import protobuf from "protobufjs";

/**
 * Register custom message types for Reward Chain
 * Uses protobufjs to properly encode/decode messages
 */
async function createRewardChainRegistry() {
  const registry = new Registry(defaultRegistryTypes);
  
  // Define MsgCreatePartner message type manually based on proto definition
  class MsgCreatePartner {
    constructor(properties = {}) {
      this.creator = properties.creator || "";
      this.name = properties.name || "";
      this.category = properties.category || "";
      this.country = properties.country || "";
      this.currency = properties.currency || "";
      this.earnCostPerPoint = properties.earnCostPerPoint || "";
      this.burnCostPerPoint = properties.burnCostPerPoint || "";
      this.totalLiquidity = properties.totalLiquidity || "";
    }
  }

  const MsgCreatePartnerType = {
    create: (properties) => new MsgCreatePartner(properties),
    encode: (message) => {
      const writer = protobuf.Writer.create();
      if (message.creator !== undefined && message.creator !== "") {
        writer.uint32(10).string(message.creator);
      }
      if (message.name !== undefined && message.name !== "") {
        writer.uint32(18).string(message.name);
      }
      if (message.category !== undefined && message.category !== "") {
        writer.uint32(26).string(message.category);
      }
      if (message.country !== undefined && message.country !== "") {
        writer.uint32(34).string(message.country);
      }
      if (message.currency !== undefined && message.currency !== "") {
        writer.uint32(42).string(message.currency);
      }
      if (message.earnCostPerPoint !== undefined && message.earnCostPerPoint !== "") {
        writer.uint32(50).string(message.earnCostPerPoint);
      }
      if (message.burnCostPerPoint !== undefined && message.burnCostPerPoint !== "") {
        writer.uint32(58).string(message.burnCostPerPoint);
      }
      if (message.totalLiquidity !== undefined && message.totalLiquidity !== "") {
        writer.uint32(66).string(message.totalLiquidity);
      }
      return writer;
    },
    decode: (input) => {
      const reader = protobuf.Reader.create(input);
      const message = {};
      const end = reader.len;
      while (reader.pos < end) {
        const tag = reader.uint32();
        switch (tag >>> 3) {
          case 1:
            message.creator = reader.string();
            break;
          case 2:
            message.name = reader.string();
            break;
          case 3:
            message.category = reader.string();
            break;
          case 4:
            message.country = reader.string();
            break;
          case 5:
            message.currency = reader.string();
            break;
          case 6:
            message.earnCostPerPoint = reader.string();
            break;
          case 7:
            message.burnCostPerPoint = reader.string();
            break;
          case 8:
            message.totalLiquidity = reader.string();
            break;
          default:
            reader.skipType(tag & 7);
            break;
        }
      }
      return message;
    },
    fromJSON: (object) => {
      return {
        creator: object.creator || "",
        name: object.name || "",
        category: object.category || "",
        country: object.country || "",
        currency: object.currency || "",
        earnCostPerPoint: object.earnCostPerPoint || "",
        burnCostPerPoint: object.burnCostPerPoint || "",
        totalLiquidity: object.totalLiquidity || "",
      };
    },
    toJSON: (message) => {
      const obj = {};
      message.creator !== undefined && (obj.creator = message.creator);
      message.name !== undefined && (obj.name = message.name);
      message.category !== undefined && (obj.category = message.category);
      message.country !== undefined && (obj.country = message.country);
      message.currency !== undefined && (obj.currency = message.currency);
      message.earnCostPerPoint !== undefined && (obj.earnCostPerPoint = message.earnCostPerPoint);
      message.burnCostPerPoint !== undefined && (obj.burnCostPerPoint = message.burnCostPerPoint);
      message.totalLiquidity !== undefined && (obj.totalLiquidity = message.totalLiquidity);
      return obj;
    },
  };
  
  registry.register("/rewardchain.rewardchain.MsgCreatePartner", MsgCreatePartnerType);

  // Define MsgAddPartnerLiquidity message type
  class MsgAddPartnerLiquidity {
    constructor(properties = {}) {
      this.creator = properties.creator || "";
      this.partnerId = properties.partnerId || 0;
      this.amount = properties.amount || "";
      this.currency = properties.currency || "";
      this.extWallet = properties.extWallet || "";
    }
  }

  const MsgAddPartnerLiquidityType = {
    create: (properties) => new MsgAddPartnerLiquidity(properties),
    encode: (message) => {
      const writer = protobuf.Writer.create();
      if (message.creator !== undefined && message.creator !== "") {
        writer.uint32(10).string(message.creator);
      }
      if (message.partnerId !== undefined && message.partnerId !== 0) {
        writer.uint32(16).uint64(Number(message.partnerId));
      }
      if (message.amount !== undefined && message.amount !== "") {
        writer.uint32(26).string(message.amount);
      }
      if (message.currency !== undefined && message.currency !== "") {
        writer.uint32(34).string(message.currency);
      }
      if (message.extWallet !== undefined && message.extWallet !== "") {
        writer.uint32(42).string(message.extWallet);
      }
      return writer;
    },
    decode: (input) => {
      const reader = protobuf.Reader.create(input);
      const message = {};
      const end = reader.len;
      while (reader.pos < end) {
        const tag = reader.uint32();
        switch (tag >>> 3) {
          case 1:
            message.creator = reader.string();
            break;
          case 2:
            message.partnerId = reader.uint64();
            break;
          case 3:
            message.amount = reader.string();
            break;
          case 4:
            message.currency = reader.string();
            break;
          case 5:
            message.extWallet = reader.string();
            break;
          default:
            reader.skipType(tag & 7);
            break;
        }
      }
      return message;
    },
    fromJSON: (object) => {
      return {
        creator: object.creator || "",
        partnerId: object.partnerId || 0,
        amount: object.amount || "",
        currency: object.currency || "",
        extWallet: object.extWallet || "",
      };
    },
    toJSON: (message) => {
      const obj = {};
      message.creator !== undefined && (obj.creator = message.creator);
      message.partnerId !== undefined && (obj.partnerId = message.partnerId);
      message.amount !== undefined && (obj.amount = message.amount);
      message.currency !== undefined && (obj.currency = message.currency);
      message.extWallet !== undefined && (obj.extWallet = message.extWallet);
      return obj;
    },
  };

  // Define MsgSwap message type
  class MsgSwap {
    constructor(properties = {}) {
      this.creator = properties.creator || "";
      this.partnerId = properties.partnerId || 0;
      this.route = properties.route || "";
      this.points = properties.points || "";
    }
  }

  const MsgSwapType = {
    create: (properties) => new MsgSwap(properties),
    encode: (message) => {
      const writer = protobuf.Writer.create();
      if (message.creator !== undefined && message.creator !== "") {
        writer.uint32(10).string(message.creator);
      }
      if (message.partnerId !== undefined && message.partnerId !== 0) {
        writer.uint32(16).uint64(Number(message.partnerId));
      }
      if (message.route !== undefined && message.route !== "") {
        writer.uint32(26).string(message.route);
      }
      if (message.points !== undefined && message.points !== "") {
        writer.uint32(34).string(message.points);
      }
      return writer;
    },
    decode: (input) => {
      const reader = protobuf.Reader.create(input);
      const message = {};
      const end = reader.len;
      while (reader.pos < end) {
        const tag = reader.uint32();
        switch (tag >>> 3) {
          case 1:
            message.creator = reader.string();
            break;
          case 2:
            message.partnerId = reader.uint64();
            break;
          case 3:
            message.route = reader.string();
            break;
          case 4:
            message.points = reader.string();
            break;
          default:
            reader.skipType(tag & 7);
            break;
        }
      }
      return message;
    },
    fromJSON: (object) => {
      return {
        creator: object.creator || "",
        partnerId: object.partnerId || 0,
        route: object.route || "",
        points: object.points || "",
      };
    },
    toJSON: (message) => {
      const obj = {};
      message.creator !== undefined && (obj.creator = message.creator);
      message.partnerId !== undefined && (obj.partnerId = message.partnerId);
      message.route !== undefined && (obj.route = message.route);
      message.points !== undefined && (obj.points = message.points);
      return obj;
    },
  };

  registry.register("/rewardchain.rewardchain.MsgAddPartnerLiquidity", MsgAddPartnerLiquidityType);
  registry.register("/rewardchain.rewardchain.MsgSwap", MsgSwapType);
  
  return registry;
}

/**
 * Reward Chain Client
 * A clean and slim client for interacting with the Reward Chain Cosmos SDK
 */
export class RewardChainClient {
  constructor(signingClient, address, rpcEndpoint, gasPrice) {
    this.client = signingClient;
    this.address = address;
    this.rpcEndpoint = rpcEndpoint;
    this.gasPrice = gasPrice;
  }

  /**
   * Create a new client instance
   * @param {string} rpcEndpoint - RPC endpoint URL (e.g., "http://localhost:26657")
   * @param {string} mnemonic - Mnemonic phrase for the wallet
   * @param {string} prefix - Bech32 address prefix (default: "reward")
   * @param {string} gasPrice - Gas price (default: "0.0001stake")
   * @returns {Promise<RewardChainClient>}
   */
  static async connect(rpcEndpoint, mnemonic, prefix = "reward", gasPrice = "0.0001stake") {
    const wallet = await DirectSecp256k1HdWallet.fromMnemonic(mnemonic, {
      prefix: prefix,
      hdPaths: [stringToPath("m/44'/118'/0'/0/0")],
    });

    const [account] = await wallet.getAccounts();
    const gasPriceObj = GasPrice.fromString(gasPrice);
    
    // Create custom registry with Reward Chain message types
    const registry = await createRewardChainRegistry();
    
    const client = await SigningStargateClient.connectWithSigner(
      rpcEndpoint,
      wallet,
      {
        gasPrice: gasPriceObj,
        registry: registry,
      }
    );

    return new RewardChainClient(client, account.address, rpcEndpoint, gasPriceObj);
  }


  /**
   * Create a new partner
   * @param {Object} partnerData - Partner data
   * @param {string} partnerData.name - Partner name
   * @param {string} partnerData.category - Partner category
   * @param {string} partnerData.country - Partner country
   * @param {string} partnerData.currency - Partner currency
   * @param {string} partnerData.earnCostPerPoint - Earn cost per point
   * @param {string} partnerData.burnCostPerPoint - Burn cost per point
   * @param {string} partnerData.totalLiquidity - Total liquidity
   * @param {Object} options - Transaction options
   * @param {string} options.memo - Transaction memo
   * @param {string} options.fee - Transaction fee (default: "auto")
   * @returns {Promise<Object>} Transaction result with partner ID
   */
  async createPartner(partnerData, options = {}) {
    const {
      name,
      category,
      country,
      currency,
      earnCostPerPoint,
      burnCostPerPoint,
      totalLiquidity,
    } = partnerData;

    const msg = {
      typeUrl: "/rewardchain.rewardchain.MsgCreatePartner",
      value: {
        creator: this.address,
        name: name,
        category: category,
        country: country,
        currency: currency,
        earnCostPerPoint: earnCostPerPoint,
        burnCostPerPoint: burnCostPerPoint,
        totalLiquidity: totalLiquidity,
      },
    };

    // Use provided fee or calculate from gas price
    let fee = options.fee;
    if (!fee || fee === "auto") {
      // Default gas limit for create partner transaction
      const gasLimit = options.gas || 200000;
      fee = calculateFee(gasLimit, this.gasPrice);
    }
    
    const memo = options.memo || "";

    const result = await this.client.signAndBroadcast(
      this.address,
      [msg],
      fee,
      memo
    );

    // Parse the response to get the partner ID
    // The response will be in the events/logs
    let partnerId = null;
    if (result.events) {
      for (const event of result.events) {
        if (event.type === "message" || event.type.includes("CreatePartner")) {
          const attributes = event.attributes || [];
          for (const attr of attributes) {
            if (attr.key === "id" || attr.key === "partner_id") {
              partnerId = attr.value;
              break;
            }
          }
        }
      }
    }

    return {
      transactionHash: result.transactionHash,
      partnerId: partnerId,
      height: result.height,
      gasUsed: result.gasUsed,
    };
  }

  /**
   * List all partners
   * @param {Object} options - Query options
   * @param {boolean} options.includeDisabled - Include disabled partners (default: false)
   * @param {Object} options.pagination - Pagination options
   * @returns {Promise<Array>} Array of partners
   */
  async listPartners(options = {}) {
    const { includeDisabled = false, pagination } = options;

    // Use REST endpoint (port 1317) for queries
    const restUrl = this.rpcEndpoint.replace(":26657", ":1317");
    const requestUrl = `${restUrl}/rewardchain/rewardchain/partners?include_disabled=${includeDisabled}`;

    try {
      const response = await fetch(requestUrl);
      if (!response.ok) {
        throw new Error(`Query failed: ${response.statusText}`);
      }
      const data = await response.json();
      return data.partners || [];
    } catch (error) {
      throw new Error(`Failed to query partners: ${error.message}`);
    }
  }

  /**
   * Get a single partner by ID
   * @param {number} partnerId - Partner ID
   * @returns {Promise<Object>} Partner data
   */
  async getPartner(partnerId) {
    const restUrl = this.rpcEndpoint.replace(":26657", ":1317");
    const requestUrl = `${restUrl}/rewardchain/rewardchain/partners/${partnerId}`;

    try {
      const response = await fetch(requestUrl);
      if (!response.ok) {
        throw new Error(`Query failed: ${response.statusText}`);
      }
      const data = await response.json();
      return data.partner || null;
    } catch (error) {
      throw new Error(`Failed to query partner: ${error.message}`);
    }
  }

  /**
   * Add liquidity for a partner
   * @param {Object} liquidityData - Liquidity data
   * @param {number} liquidityData.partnerId - Partner ID
   * @param {string} liquidityData.amount - Amount to add
   * @param {string} liquidityData.currency - Currency
   * @param {string} liquidityData.extWallet - External wallet address
   * @param {Object} options - Transaction options
   * @param {string} options.memo - Transaction memo
   * @param {string|Object} options.fee - Transaction fee (default: calculated)
   * @param {number} options.gas - Gas limit (default: 200000)
   * @returns {Promise<Object>} Transaction result
   */
  async addPartnerLiquidity(liquidityData, options = {}) {
    const {
      partnerId,
      amount,
      currency,
      extWallet,
    } = liquidityData;

    const msg = {
      typeUrl: "/rewardchain.rewardchain.MsgAddPartnerLiquidity",
      value: {
        creator: this.address,
        partnerId: partnerId,
        amount: amount,
        currency: currency,
        extWallet: extWallet,
      },
    };

    // Use provided fee or calculate from gas price
    let fee = options.fee;
    if (!fee || fee === "auto") {
      const gasLimit = options.gas || 200000;
      fee = calculateFee(gasLimit, this.gasPrice);
    }
    
    const memo = options.memo || "";

    const result = await this.client.signAndBroadcast(
      this.address,
      [msg],
      fee,
      memo
    );

    return {
      transactionHash: result.transactionHash,
      height: result.height,
      gasUsed: result.gasUsed,
    };
  }

  /**
   * Swap between points and tokens for a partner
   * @param {Object} swapData - Swap data
   * @param {number} swapData.partnerId - Partner ID
   * @param {string} swapData.route - Swap route: "points_to_token" or "token_to_points"
   * @param {string} swapData.points - Points amount
   * @param {Object} options - Transaction options
   * @param {string} options.memo - Transaction memo
   * @param {string|Object} options.fee - Transaction fee (default: calculated)
   * @param {number} options.gas - Gas limit (default: 200000)
   * @returns {Promise<Object>} Transaction result
   */
  async swap(swapData, options = {}) {
    const {
      partnerId,
      route,
      points,
    } = swapData;

    // Validate route
    if (route !== "points_to_token" && route !== "token_to_points") {
      throw new Error('Route must be either "points_to_token" or "token_to_points"');
    }

    const msg = {
      typeUrl: "/rewardchain.rewardchain.MsgSwap",
      value: {
        creator: this.address,
        partnerId: partnerId,
        route: route,
        points: points,
      },
    };

    // Use provided fee or calculate from gas price
    let fee = options.fee;
    if (!fee || fee === "auto") {
      const gasLimit = options.gas || 200000;
      fee = calculateFee(gasLimit, this.gasPrice);
    }
    
    const memo = options.memo || "";

    const result = await this.client.signAndBroadcast(
      this.address,
      [msg],
      fee,
      memo
    );

    return {
      transactionHash: result.transactionHash,
      height: result.height,
      gasUsed: result.gasUsed,
    };
  }

  /**
   * Disconnect the client
   */
  async disconnect() {
    await this.client.disconnect();
  }
}

export default RewardChainClient;
