import RewardChainClient from "./client.js";

/**
 * Example usage of the Reward Chain Client
 * 
 * Make sure to:
 * 1. Set your RPC endpoint
 * 2. Set your mnemonic phrase
 * 3. Adjust the address prefix if needed
 */

// Configuration
const RPC_ENDPOINT = process.env.RPC_ENDPOINT || "http://localhost:26657";
const MNEMONIC = process.env.MNEMONIC || "vote between double monitor month circle before hybrid praise hotel brown rally skill sunny buffalo foot region define pledge adult medal whip soft account";
const ADDRESS_PREFIX = process.env.ADDRESS_PREFIX || "reward";

async function main() {
  try {
    console.log("Connecting to Reward Chain...");
    const client = await RewardChainClient.connect(
      RPC_ENDPOINT,
      MNEMONIC,
      ADDRESS_PREFIX
    );
    console.log(`Connected! Address: ${client.address}\n`);

    // Example 1: Create a partner
    console.log("Creating a new partner...");
    const createResult = await client.createPartner({
      name: "Example Partner",
      category: "Retail",
      country: "US",
      currency: "USD",
      earnCostPerPoint: "1.0",
      burnCostPerPoint: "0.9",
      totalLiquidity: "10000",
    });

    console.log("Partner created successfully!");
    console.log("Transaction Hash:", createResult.transactionHash);
    console.log("Partner ID:", createResult.partnerId);
    console.log("Height:", createResult.height);
    console.log("Gas Used:", createResult.gasUsed);
    console.log("\n");

    // Wait a bit for the transaction to be indexed
    await new Promise((resolve) => setTimeout(resolve, 2000));

    // Example 2: List all partners
    console.log("Listing all partners...");
    const partners = await client.listPartners({
      includeDisabled: false,
    });

    console.log(`Found ${partners.length} partner(s):`);
    partners.forEach((partner, index) => {
      console.log(`\nPartner ${index + 1}:`);
      console.log("  ID:", partner.id);
      console.log("  Name:", partner.name);
      console.log("  Category:", partner.category);
      console.log("  Country:", partner.country);
      console.log("  Total Liquidity:", partner.total_liquidity);
      console.log("  Available Liquidity:", partner.available_liquidity);
      console.log("  Disabled:", partner.disabled);
    });

    // Example 3: Get a specific partner (if we have an ID)
    if (partners.length > 0) {
      const partnerId = partners[0].id;
      console.log(`\nFetching partner with ID ${partnerId}...`);
      const partner = await client.getPartner(partnerId);
      console.log("Partner details:", JSON.stringify(partner, null, 2));

      // Example 4: Add liquidity to a partner
      console.log(`\nAdding liquidity to partner ${partnerId}...`);
      const liquidityResult = await client.addPartnerLiquidity({
        partnerId: partnerId,
        amount: "1000",
        currency: "USD",
        extWallet: "0x1234567890123456789012345678901234567890",
      });
      console.log("Liquidity added successfully!");
      console.log("Transaction Hash:", liquidityResult.transactionHash);
      console.log("Height:", liquidityResult.height);
      console.log("Gas Used:", liquidityResult.gasUsed);

      // Wait a bit for the transaction to be indexed
      await new Promise((resolve) => setTimeout(resolve, 2000));

      // Example 5: Swap points to tokens
      console.log(`\nSwapping points to tokens for partner ${partnerId}...`);
      const swapResult = await client.swap({
        partnerId: partnerId,
        route: "points_to_token",
        points: "100",
      });
      console.log("Swap completed successfully!");
      console.log("Transaction Hash:", swapResult.transactionHash);
      console.log("Height:", swapResult.height);
      console.log("Gas Used:", swapResult.gasUsed);
    }

    // Disconnect
    await client.disconnect();
    console.log("\nDisconnected from Reward Chain");
  } catch (error) {
    console.error("Error:", error.message);
    if (error.stack) {
      console.error(error.stack);
    }
    process.exit(1);
  }
}

// Run the example
main();
